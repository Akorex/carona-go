package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	// Enable strict JSON decoding: reject any unexpected/unknown fields in request bodies
	binding.EnableDecoderDisallowUnknownFields = true

	// Tell Gin's validator to use the JSON tag name instead of the Go struct field name
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

// FieldError represents a single validation failure on a field.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormatValidationErrors(err error) []FieldError {
	// Handle strict mode unknown field error (e.g. 'json: unknown field "age"')
	if strings.HasPrefix(err.Error(), "json: unknown field ") {
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		fieldName = strings.Trim(fieldName, "\"")
		return []FieldError{
			{Field: fieldName, Message: "unknown or unexpected field"},
		}
	}

	var valErrors validator.ValidationErrors
	if !errors.As(err, &valErrors) {
		return []FieldError{
			{Field: "body", Message: "Invalid JSON format"},
		}
	}

	var errorsList []FieldError
	for _, fieldErr := range valErrors {
		errorsList = append(errorsList, FieldError{
			Field:   fieldErr.Field(),
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
