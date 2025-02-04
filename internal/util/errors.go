package util

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var JsonDecodeError = errors.New("Error while decoding json")

func FieldErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field required"
	case "email":
		return "Invalid email"
	default:
		return fe.Error()
	}
}

func GetErrorMessages(ve validator.ValidationErrors) string {
	result := ""
	for _, err := range ve {
		result += err.Field() + ": " + FieldErrorMessage(err) + "\n"
	}
	return result
}
