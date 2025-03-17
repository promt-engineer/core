package validator

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/constants"
	"fmt"
	"sync"
)

func (v *Validator) generateErrorMessage() {
	errorMessagesOnce.Do(func() {
		errorMessages = map[string]string{
			"required":      "field %s is required",
			"email_custom":  "email %s is not valid",
			"str_gt":        "field %s must have greater than %s characters",
			"str_lt":        "field %s must have less than %s characters",
			"has_lowercase": "field %s must have at least one small character",
			"has_uppercase": "field %s must have at least one big character",
			"has_special":   "field %s must have at least one special character",
			"oneof":         "field %s must have value one of allowed list: %s",
			"gte":           "field %s must be greater or equal than %s",
			"gt":            "field %s must be greater than %s",
			"lte":           "field %s must be less or equal than %s",
			"lt":            "field %s must be less than %s",
			"url":           "field %s must be an url",
			"uuid":          "field %s must be an uuid",

			constants.GameRuleName: "field %s must be one of " + fmt.Sprintf("%v", v.config.AvailableGames),
		}
	})
}

var (
	errorMessages     map[string]string
	errorMessagesOnce sync.Once
)
