package middleware

import (
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

//	"github.com/0987363/models"

	"net"
	"net/url"
	"os"
	"time"
)

const loggerKey = "Logger"

var logConn net.Conn

func LoggerConnInit() {
	if viper.GetString("log.dst") != "" {
		conn, err := net.Dial("udp", viper.GetString("log.dst"))
		if err != nil {
			logrus.Fatal(err)
		}
		logConn = conn
	}
}

func LoggerInit() *logrus.Logger {
	logger := logrus.New()
//	host, _ := os.Hostname()

	logger.Out = os.Stdout
//	logger.Level = models.ConvertLevel(viper.GetString("log.level"))
	logger.Formatter = &logrus.TextFormatter{FullTimestamp: true, TimestampFormat: time.RFC3339Nano}

	if logConn == nil {
		return logger
	}

	/*
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
	*/

	return logger
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := GetLogger(c)
		if log == nil {
			log = LoggerInit().WithField(RequestIDKey, GetRequestID(c))
		}

		u, err := url.QueryUnescape(c.Request.URL.String())
		if err != nil {
			u = c.Request.URL.String()
		}
		log = log.WithFields(
			logrus.Fields{
				"method":     c.Request.Method,
				"user_agent": c.Request.UserAgent(),
				"url":        u,
				"remote":     c.ClientIP(),
			})

		start := time.Now()

		c.Set(loggerKey, log)
		c.Next()

		log.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"user_agent": c.Request.UserAgent(),
			"url":        u,
			"status":     c.Writer.Status(),
			"remote":     c.ClientIP(),
			"spend":      time.Now().Sub(start).String(),
		}).Infof("Responded %03d in %s", c.Writer.Status(), time.Now().Sub(start))
	}
}

func GetLogger(c *gin.Context) *logrus.Entry {
	if logger, ok := c.Get(loggerKey); ok {
		return logger.(*logrus.Entry)
	}

	return nil
}
