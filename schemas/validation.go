package schemas

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New()

func ValidateRequest(r interface{}) ([]error, error) {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	issues := make([]error, 0)
	err := validate.Struct(r)
	if err != nil {
		for _, field := range err.(validator.ValidationErrors) {
			issues = append(issues, field)
		}
	}
	return issues, err
}
