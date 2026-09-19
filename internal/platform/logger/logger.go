package logger

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

var BaseLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelInfo,
}))

func FromContext(c *gin.Context) *slog.Logger {
	if reqID, exists := c.Get("requestId"); exists {
		return BaseLogger.With("request_id", reqID)
	}

	return BaseLogger
}
