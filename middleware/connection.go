package middleware

import (
	"net"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const PushConnectionKey = "PushConnection"

var pushConn *net.Conn

func ConnectionInit() {
	if address := viper.GetString("push_address"); address != "" {
		conn, err := net.Dial("udp", address)
		if err == nil {
			pushConn = &conn
			log.Info("Dial push connection success:", address)
		} else {
			log.Info("Dial push connection failed:", err)
		}
		return
	}
	log.Info("Did not set push address.")
}

func PushConnection() gin.HandlerFunc {
	return func(c *gin.Context) {
		if pushConn == nil {
			log.Info("Start init push connection success.")
			ConnectionInit()
		}

		conn := GetPushConnection(c)
		if conn != nil {
			c.Set(PushConnectionKey, conn)
		}

		c.Next()
	}
}

func GetPushConnection(c *gin.Context) *net.Conn {
	if conn, ok := c.Get(PushConnectionKey); ok {
		return conn.(*net.Conn)
	}

	return nil
}
