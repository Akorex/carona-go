package health

import (
	"caronago/internal/platform/response"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Check(c *gin.Context) {
	response.SendOK(c, "Carona API is running smoothly", gin.H{
		"status": "healthy",
	})
}

func RegisterRoutes(router *gin.Engine) {
	h := NewHandler()
	router.GET("/health", h.Check)
}
