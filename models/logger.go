package models

import (
	"net"
	"os"
	"time"

	"github.com/0987363/viper"
	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
)

func LoggerInit(method string) *logrus.Entry {
	logger := logrus.New()
	logger.Level = ConvertLevel(viper.GetString("log.level"))
	logger.Formatter = &logrus.TextFormatter{FullTimestamp: true, TimestampFormat: time.RFC3339Nano}
	logger.Out = os.Stdout

	if dst := viper.GetString("log.dst"); dst != "" {
		logConn, err := net.Dial("udp", dst)
		if err != nil {
			logger.Fatal("Dial remote logrustash failed: ", err)
		}

		host, _ := os.Hostname()
		hook := logrustash.New(logConn, logrustash.LogstashFormatter{
			Fields: logrus.Fields{
				"type":     "platform",
				"hostname": host,
				"service":  "configGO",
				"release":  viper.GetString("release"),
			},
			Formatter: &logrus.JSONFormatter{
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyTime: "@timestamp",
					logrus.FieldKeyMsg:  "message",
				},
				TimestampFormat: time.RFC3339Nano,
			}})
		logger.Hooks.Add(hook)
	}

	return logger.WithFields(
		logrus.Fields{
			"method": method,
		})
}

func ConvertLevel(level string) logrus.Level {
	switch level {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}
