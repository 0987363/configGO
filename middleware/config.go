package middleware

import (
	"github.com/gin-gonic/gin"
)

const configKey = "ConfigGO"

func Config(m interface{}, handle gin.HandlerFunc) gin.HandlerFunc {
	if m == nil {
		return handle
	}
	return func(c *gin.Context) {
		c.Set(configKey, m)
		handle(c)
	}
}

func GetConfig(c *gin.Context) interface{} {
	if cg, ok := c.Get(configKey); ok {
		return cg
	}

	return nil
}
