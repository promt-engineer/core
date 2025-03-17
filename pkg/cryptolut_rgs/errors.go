package cryptolut_rgs

import (
	"errors"
)

var ErrCLSessionNotFound = errors.New("session not found")
var ErrCLNotConfigured = errors.New("not configured")
var ErrCLAccessDenied = errors.New("access denied")
var ErrCLValidationError = errors.New("validation error")
var ErrCLInternalServerError = errors.New("internal server error")
var ErrCLCurrencyNotFound = errors.New("currency not found")
var ErrCLCountryNotFound = errors.New("country not found")
var ErrCLNotAvailable = errors.New("not available")
var ErrCLSessionTimeout = errors.New("session timeout")
var ErrCLSessionDuplicate = errors.New("session duplicated")
var ErrCLSessionReopened = errors.New("session reopened")

var errorMap = map[string]error{
	"SESSION_NOT_FOUND":     ErrCLSessionNotFound,
	"NOT_CONFIGURED":        ErrCLNotConfigured,
	"ACCESS_DENIED":         ErrCLAccessDenied,
	"VALIDATION_ERROR":      ErrCLValidationError,
	"INTERNAL_SERVER_ERROR": ErrCLInternalServerError,
	"CURRENCY_NOT_FOUND":    ErrCLCurrencyNotFound,
	"COUNTRY_NOT_FOUND":     ErrCLCountryNotFound,
	"NOT_AVAILABLE":         ErrCLNotAvailable,
	"SESSION_TIMEOUT":       ErrCLSessionTimeout,
	"SESSION_DUPLICATE":     ErrCLSessionDuplicate,
	"SESSION_REOPENED":      ErrCLSessionReopened,
}
