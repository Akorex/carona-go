package middleware

import (
	"caronago/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// process request
		c.Next()

		latency := time.Since(start).Milliseconds()
		status := c.Writer.Status()
		log := logger.FromContext(c)

		attributes := []any{
			"status", status,
			"method", method,
			"path", path,
			"latency_ms", latency,
			"client_ip", c.ClientIP(),
		}

		if status >= 500 {
			log.Error("server error", attributes...)
		} else if status >= 400 {
			log.Warn("client error", attributes...)
		} else {
			log.Info("request completed", attributes...)
		}

	}
}
