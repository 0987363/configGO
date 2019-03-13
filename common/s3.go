package common

import (
	"context"
	"strings"
	"time"

	"github.com/0987363/configGO/models"
	"github.com/0987363/viper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/radovskyb/watcher"
)

func WatchS3() {
	log := models.LoggerInit("watch_s3")

	if viper.GetString("aws.s3.key") == "" {
		log.Fatal("Aws s3 config invalid.")
	}

	sess := session.New(&aws.Config{
		Region:      aws.String(viper.GetString("aws.sqs.region")),
		MaxRetries:  aws.Int(5),
		Credentials: credentials.NewStaticCredentials(viper.GetString("aws.s3.key"), viper.GetString("aws.s3.secret"), ""),
	})
	svc := s3.New(sess)

	go func() {
		for {
			select {
			case notify := <-chNotify:
				if err := notify.Init(); err != nil {
					log.Info("Init notify failed: ", err)
					continue
				}
				s := NewService(notify.Key)
				if s == nil {
					notify.Ch <- notify
					continue
				}
				go func() {
					defer func() {
						notify.Ch <- notify
					}()

					events := strings.Split(notify.Records[0].EventName, ":")
					if len(events) < 2 {
						return
					}
					log.Info("Recv event: ", events)
					switch events[0] {
					case "ObjectCreated":
						ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
						defer cancel()

						res, err := svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
							Bucket: aws.String(notify.Bucket),
							Key:    aws.String(notify.Key),
						})
						if err != nil {
							if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
								log.Error("download timeout, ", err)
							} else {
								log.Error("failed to download object, ", err, notify.Bucket, notify.Key)
							}
							return
						}
						s.Load(res.Body)
						break
					case "ObjectRemoved":
						s.Op = watcher.Remove
						break
					default:
						log.Warning("Recv unknown event: ", events)
						return
					}

					log.Info("Send json message to etcd watch: ", s.Key())
					chService <- s
				}()
			}
		}
	}()
}
