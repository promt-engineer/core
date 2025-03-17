package errs

import "errors"

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
