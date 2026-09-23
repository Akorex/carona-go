package auth

import (
	"caronago/internal/platform/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	service := NewService(db, cfg.JWTSecret, cfg.JWTExpiresInHours)
	handler := NewHandler(service)

	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", handler.Register)
	}
}
