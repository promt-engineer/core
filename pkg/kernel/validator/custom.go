package validator

import (
	"github.com/go-playground/validator/v10"
)

func stringOneOfTheList(list []string) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		str, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}

		for _, game := range list {
			if str == game {
				return true
			}
		}

		return false
	}
}

func stringOneOfTheListOrNil(list []string) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		if fl.Field().Interface() == nil {
			return true
		}

		str, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}

		for _, game := range list {
			if str == game {
				return true
			}
		}

		return false
	}
}