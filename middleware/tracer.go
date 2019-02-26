package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Tracer records the start and end of each request through logger
func Tracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := GetLogger(c)

		start := time.Now()
		c.Next()
		end := time.Now()

		endTrace(logger, c, end.Sub(start))
	}
}

func startTrace(logger *logrus.Entry, r *http.Request) {
	if logger == nil {
		return
	}

	logger.Infof("%s %q from %s", r.Method, r.URL.String(), r.RemoteAddr)
}

func endTrace(logger *logrus.Entry, c *gin.Context, d time.Duration) {
	if logger == nil {
		return
	}

	logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url":    c.Request.URL.String(),
		"status": c.Writer.Status(),
		"spend":  d.String(),
	}).Infof("Responded %s in %s", c.Writer.Status(), d)
}
