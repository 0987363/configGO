package middleware

import (
	//	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recoverer catches the panic from the handler gracefully and logs the error
// with stack trace.
func Recoverer() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := GetLogger(c)

		defer func() {
			if err := recover(); err != nil {
				if logger != nil {
					logger.Errorf("recv stack: %s\n%s", err, string(debug.Stack()))
				}

				//	writeInternalServerError(c)
				c.Abort()
				c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}
		}()
		c.Next()
	}
}

func writeInternalServerError(c *gin.Context) {
	c.Abort()

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("X-Content-Type-Options", "nosniff")

	c.String(http.StatusInternalServerError, "InternalServerError")
}
