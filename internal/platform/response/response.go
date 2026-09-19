package response

import (
	"net/http"

	"caronago/internal/platform/apierror"

	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func SendSuccess(c *gin.Context, code int, message string, data any) {
	c.JSON(code, Envelope{
		Success: true,
		Message: message,
		Code:    code,
		Data:    data,
	})
}

func SendError(c *gin.Context, code int, message string, err any) {
	c.JSON(code, Envelope{
		Success: false,
		Message: message,
		Code:    code,
		Data:    err,
	})
}

func SendOK(c *gin.Context, message string, data any) {
	SendSuccess(c, http.StatusOK, message, data)
}

func SendCreated(c *gin.Context, message string, data any) {
	SendSuccess(c, http.StatusCreated, message, data)
}

func FromApiError(c *gin.Context, err error) {
	code, msg := apierror.FromError(err)

	SendError(c, code, msg, nil)
}
