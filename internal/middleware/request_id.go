package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"
const RequestIDContextKey = "requestId"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDHeader)
		if id == "" {
			id = uuid.New().String()
		}

		// Store in Gin context for handlers to access
		c.Set(RequestIDContextKey, id)

		c.Header(RequestIDHeader, id)

		c.Next()
	}
}
