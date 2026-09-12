package middleware

import (
	"caronago/internal/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NoRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		msg := fmt.Sprintf("Route %s %s not found", c.Request.Method, c.Request.URL.Path)
		response.SendError(c, http.StatusNotFound, msg, nil)
	}
}

func NoMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		msg := fmt.Sprintf("Method %s not allowed on %s", c.Request.Method, c.Request.URL.Path)
		response.SendError(c, http.StatusMethodNotAllowed, msg, nil)
	}
}
