package helpers

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func ValidationErrorFormatter(err error, input interface{}) map[string]string {
	validationErrors := err.(validator.ValidationErrors)
	errorFields := make(map[string]string)

	reflected := reflect.TypeOf(input)
	for _, err := range validationErrors {
		field, _ := reflected.FieldByName(err.Field())
		jsonTag := field.Tag.Get("json")

		switch err.Tag() {
		case "required":
			errorFields[jsonTag] = "This field is required"
		case "email":
			errorFields[jsonTag] = "Invalid email format"
		case "min":
			if err.Field() == "Password" {
				errorFields[jsonTag] = "Password must be at least 8 characters"
			}
		case "eqfield":
			errorFields[jsonTag] = "Password Confirm do not match"
		case "e164":
			errorFields[jsonTag] = "Invalid phone number format"
		}
	}

	return errorFields
}
