package auth

import (
	"caronago/internal/platform/response"
	"caronago/internal/platform/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		fieldErrors := validator.FormatValidationErrors(err)
		response.SendError(c, http.StatusUnprocessableEntity, "Validation failed", fieldErrors)
		return
	}

	res, err := h.service.Register(c.Request.Context(), input)
	if err != nil {
		response.FromApiError(c, err)
		return
	}

	response.SendCreated(c, "User registered successfully", res)

}
