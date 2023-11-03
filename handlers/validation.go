package handlers

import (
	"github.com/go-playground/validator/v10"
)

func ValidateRequest[Request any](r Request) error {

	var validate = validator.New()
	return validate.Struct(r)
}
