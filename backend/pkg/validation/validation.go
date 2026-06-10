// Package validation translates go-playground validator errors into a flat
// field -> message map for the API error envelope.
package validation

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// FieldErrors converts a validator.ValidationErrors into a field->message map.
// Non-validation errors return nil so callers can fall back to a generic message.
func FieldErrors(err error) map[string]string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return nil
	}
	out := make(map[string]string, len(ve))
	for _, fe := range ve {
		out[fe.Field()] = message(fe)
	}
	return out
}

func message(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return fmt.Sprintf("Must be at least %s characters/value", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters/value", fe.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fe.Param())
	case "e164":
		return "Must be a valid phone number"
	default:
		return fmt.Sprintf("Failed validation: %s", fe.Tag())
	}
}
