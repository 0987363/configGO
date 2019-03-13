package common

import (
	"encoding/json"
	"time"

	"github.com/0987363/configGO/models"
	"github.com/0987363/viper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type AwsSqs struct {
	Records []struct {
		EventVersion string    `json:"eventVersion"`
		EventSource  string    `json:"eventSource"`
		AwsRegion    string    `json:"awsRegion"`
		EventTime    time.Time `json:"eventTime"`
		EventName    string    `json:"eventName"`
		UserIdentity struct {
			PrincipalID string `json:"principalId"`
		} `json:"userIdentity"`
		RequestParameters struct {
			SourceIPAddress string `json:"sourceIPAddress"`
		} `json:"requestParameters"`
		ResponseElements struct {
			XAmzRequestID string `json:"x-amz-request-id"`
			XAmzID2       string `json:"x-amz-id-2"`
		} `json:"responseElements"`
		S3 struct {
			S3SchemaVersion string `json:"s3SchemaVersion"`
			ConfigurationID string `json:"configurationId"`
			Bucket          struct {
				Name          string `json:"name"`
				OwnerIdentity struct {
					PrincipalID string `json:"principalId"`
				} `json:"ownerIdentity"`
				Arn string `json:"arn"`
			} `json:"bucket"`
			Object struct {
				Key       string `json:"key"`
				Size      int    `json:"size"`
				ETag      string `json:"eTag"`
				VersionID string `json:"versionId"`
				Sequencer string `json:"sequencer"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

type AwsNotify struct {
	Bucket string
	Key    string
	Ch     chan *AwsNotify
	Msg    *sqs.Message
	AwsSqs
}

var chNotify chan *AwsNotify

func init() {
	chNotify = make(chan *AwsNotify, 10)
}

func WatchSqs() {
	sess := session.New(&aws.Config{
		Region:      aws.String(viper.GetString("aws.sqs.region")),
		MaxRetries:  aws.Int(5),
		Credentials: credentials.NewStaticCredentials(viper.GetString("aws.sqs.key"), viper.GetString("aws.sqs.secret"), ""),
	})

	svc := sqs.New(sess)
	receive_params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(viper.GetString("aws.sqs.queue")),
		MaxNumberOfMessages: aws.Int64(10), // 一次最多取幾個 message
		VisibilityTimeout:   aws.Int64(5),  // 如果這個 message 沒刪除，下次再被取出來的時間
		WaitTimeSeconds:     aws.Int64(20), // long polling 方式取，會建立一條長連線並且等在那邊，直到 SQS 收到新 message 回傳給這條連線才中斷
	}

	chFinish := make(chan *AwsNotify, 10)

	log := models.LoggerInit("watch_sqs")

	go func() {
		for {
			select {
			case notify := <-chFinish:
				delete_params := &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(viper.GetString("aws.sqs.queue")),
					ReceiptHandle: notify.Msg.ReceiptHandle,
				}

				log.Infof("Delete message ID: %s, bucket:%s, key:%s",
					*notify.Msg.MessageId, notify.Bucket, notify.Key)
				_, err := svc.DeleteMessage(delete_params) // No response returned when successed.
				if err != nil {
					log.Errorf("Delete message ID: %s has beed deleted.\n", *notify.Msg.MessageId)
				}
			}
		}
	}()

	for {
		resp, err := svc.ReceiveMessage(receive_params)
		if err != nil {
			log.Error("Recv sqs message failed: ", err)
			continue
		}
		for _, message := range resp.Messages {
			notify := AwsNotify{}
			if err := json.Unmarshal([]byte(*message.Body), &notify); err != nil {
				log.Error("Unmarshal message failed: ", err)
				continue
			}

			log.Infof("Send message to s3 watch, bucket:%s, key:%s.", notify.Records[0].S3.Bucket.Name, notify.Records[0].S3.Object.Key)
			notify.Bucket = notify.Records[0].S3.Bucket.Name
			notify.Key = notify.Records[0].S3.Object.Key
			notify.Ch = chFinish
			notify.Msg = message
			chNotify <- &notify
		}
	}
}
