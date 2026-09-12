package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FieldError represents a single validation failure on a field.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatValidationErrors inspects any error returned by Gin's c.ShouldBindJSON().
// If it's a validator.ValidationErrors, it transforms it into a clean slice of FieldErrors.
func FormatValidationErrors(err error) []FieldError {
	var valErrors validator.ValidationErrors
	if !errors.As(err, &valErrors) {
		return []FieldError{
			{Field: "body", Message: "Invalid JSON format"},
		}
	}

	var errorsList []FieldError
	for _, fieldErr := range valErrors {
		errorsList = append(errorsList, FieldError{
			Field:   strings.ToLower(fieldErr.Field()),
			Message: msgForTag(fieldErr),
		})
	}

	return errorsList
}

// msgForTag converts cryptic validator tags into human-readable messages.
func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters long", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters long", fe.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", fe.Param())
	default:
		return fmt.Sprintf("failed validation on '%s'", fe.Tag())
	}
}
