package validatorHandler

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func SetValidator() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func GetValidator() *validator.Validate {
	return validate
}
