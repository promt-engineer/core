package errs

import "strings"

type Action string

const (
	ActionRestart  Action = "restart"
	ActionContinue Action = "continue"
)

type GameHubError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Display bool   `json:"display"`
	Action  Action `json:"action"`
}

type GameHubResponse struct {
	Status int          `json:"status"`
	Error  GameHubError `json:"error"`
}

const (
	ErrCodeUnknown           = "ERR001"
	ErrCodeSessionExpired    = "ERR002"
	ErrCodeInsufficientFunds = "ERR003"
	ErrCodeWageringLimit     = "ERR004"
	ErrCodeAuthFailed        = "ERR005"
	ErrCodeUnauthorized      = "ERR006"
	ErrCodeDuplicate         = "ERR007"
	ErrCodeCurrency          = "ERR008"
)

var GameHubErrorMap = map[string]GameHubError{
	ErrCodeUnknown: {
		Code:    ErrCodeUnknown,
		Display: true,
		Action:  ActionRestart,
		Message: "Unknown error occurred",
	},
	ErrCodeSessionExpired: {
		Code:    ErrCodeSessionExpired,
		Display: true,
		Action:  ActionRestart,
		Message: "The session has timed out. Please login again to continue playing",
	},
	ErrCodeInsufficientFunds: {
		Code:    ErrCodeInsufficientFunds,
		Display: true,
		Action:  ActionContinue,
		Message: "Insufficient funds to place current wager. Please reduce the stake or add more funds to your balance",
	},
	ErrCodeWageringLimit: {
		Code:    ErrCodeWageringLimit,
		Display: true,
		Action:  ActionContinue,
		Message: "This wagering will exceed your wagering limitation. Please try a smaller amount or increase the limit",
	},
	ErrCodeAuthFailed: {
		Code:    ErrCodeAuthFailed,
		Display: true,
		Action:  ActionRestart,
		Message: "Player authentication failed",
	},
	ErrCodeUnauthorized: {
		Code:    ErrCodeUnauthorized,
		Display: false,
		Action:  ActionRestart,
		Message: "Unauthorized request",
	},
	ErrCodeDuplicate: {
		Code:    ErrCodeDuplicate,
		Display: true,
		Action:  ActionContinue,
		Message: "Duplicate transaction request, means this transaction was already processed",
	},
	ErrCodeCurrency: {
		Code:    ErrCodeCurrency,
		Display: true,
		Action:  ActionRestart,
		Message: "Unsupported currency",
	},
}

func GetGameHubError(code string) (GameHubError, bool) {
	err, exists := GameHubErrorMap[code]
	return err, exists
}

func MapErrorToGameHub(err error) (GameHubError, bool) {
	if err == nil {
		return GameHubError{}, false
	}

	errMsg := err.Error()

	switch err {
	case ErrBalanceTooLow, ErrNotEnoughMoney:
		return GameHubErrorMap[ErrCodeInsufficientFunds], true
	case ErrSessionTokenExpired:
		return GameHubErrorMap[ErrCodeSessionExpired], true
	case ErrWrongSessionToken:
		return GameHubErrorMap[ErrCodeUnauthorized], true
	case ErrUserIsBlocked:
		return GameHubErrorMap[ErrCodeAuthFailed], true
	case ErrUserHasDifferentCurrency:
		return GameHubErrorMap[ErrCodeCurrency], true
	}

	switch {
	case containsAny(errMsg,
		"Insufficient funds",
		"not enough money",
		"low balance"):
		return GameHubErrorMap[ErrCodeInsufficientFunds], true

	case containsAny(errMsg,
		"session has timed out",
		"session token expired",
		"session expired"):
		return GameHubErrorMap[ErrCodeSessionExpired], true

	case containsAny(errMsg,
		"authentication failed",
		"user is blocked"):
		return GameHubErrorMap[ErrCodeAuthFailed], true

	case containsAny(errMsg,
		"Unauthorized",
		"wrong session token"):
		return GameHubErrorMap[ErrCodeUnauthorized], true

	case containsAny(errMsg,
		"Duplicate transaction",
		"already processed"):
		return GameHubErrorMap[ErrCodeDuplicate], true

	case containsAny(errMsg,
		"Unsupported currency",
		"different currency"):
		return GameHubErrorMap[ErrCodeCurrency], true

	case containsAny(errMsg,
		"wagering limitation",
		"wager limit exceeded"):
		return GameHubErrorMap[ErrCodeWageringLimit], true
	}

	return GameHubErrorMap[ErrCodeUnknown], true
}

func containsAny(s string, substrs ...string) bool {
	s = strings.ToLower(s)
	for _, substr := range substrs {
		if strings.Contains(s, strings.ToLower(substr)) {
			return true
		}
	}
	return false
}
