package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// We use this to generate human readable error messages for validation errors.
func MsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s is not a valid email", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", fe.Field(), fe.Param())
	case "required_with":
		return fmt.Sprintf("%s is required", fe.Field())
	}
	return fe.Error() // default error
}
