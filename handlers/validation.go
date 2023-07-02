package handlers

import (
	"github.com/go-playground/validator/v10"
)

func ValidateRequest[Request any](r Request) ([]string, error) {

	var validate = validator.New()

	var badFields []string
	err := validate.Struct(r)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			badFields = append(badFields, err.StructNamespace())
		}
	}
	return badFields, err
}
