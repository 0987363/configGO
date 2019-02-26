package middleware

import (
	//	"net/http"

	//	"golang.org/x/net/context"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

// RequestIDHeader is the http header for Request ID
const RequestIDHeader = "X-Request-Id"

// RequestIDKey is the context key for request ID
const RequestIDKey = "RequestID"

// RequestID middleware put a request ID into the context object. If the request
// does not exist in the request header, a new one will be generated
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetRequestID(c) == "" {

			// Use the header if provided
			id := c.Request.Header.Get(RequestIDHeader)

			// Otherwise generate a new UUID
			if id == "" {
				id = uuid.NewV4().String()
			}

			c.Set(RequestIDKey, id)
			c.Next()
		}
	}
}

// GetRequestID returns request id from the context, or empty string if the
// RequestID middleware is not used.
func GetRequestID(c *gin.Context) string {
	if requestID, ok := c.Get(RequestIDKey); ok {
		return requestID.(string)
	}

	return ""
}
