package utils

import "github.com/go-playground/validator/v10"

func FormatValidationError(err error) map[string]string {
	errors := map[string]string{}

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			errors[fe.Field()] = fe.Tag()
		}
	}

	return errors
}
