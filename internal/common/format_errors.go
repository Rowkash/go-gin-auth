package common

import (
	"fmt"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormatValidationErrors(validationErrors validator.ValidationErrors) []ValidationError {
	var errors []ValidationError

	for _, e := range validationErrors {
		errors = append(errors, ValidationError{
			Field:   lowerFirst(e.Field()),
			Message: getErrorMessage(e),
		})
	}

	return errors
}

func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Must be at least %s characters", e.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters", e.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", e.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", e.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", e.Param())
	case "eqfield":
		return fmt.Sprintf("Must match %s", e.Param())
	//case "phone":
	//	return "Invalid phone number format (use E.164: +1234567890)"
	case "slug":
		return "Must be a valid URL slug (lowercase letters, numbers, hyphens)"
	default:
		return fmt.Sprintf("Failed validation: %s", e.Tag())
	}
}

func lowerFirst(str string) string {
	if str == "" {
		return ""
	}
	runes := []rune(str)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
