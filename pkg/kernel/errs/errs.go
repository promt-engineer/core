package errs

import (
	"errors"
	"fmt"
)

type InternalValidationError struct {
	Err error
}

func (err InternalValidationError) Error() string {
	return err.Err.Error()
}

func NewInternalValidationError(err error) InternalValidationError {
	return InternalValidationError{Err: err}
}

func NewInternalValidationErrorFromString(str string) InternalValidationError {
	return NewInternalValidationError(errors.New(str))
}

var (
	ErrUpdatingIsNotAllowed = errors.New("updating spin indexes is not allowed")
	ErrLastSpinWasNotShown  = errors.New("last spin was not shown")

	ErrNotEnoughMoney = errors.New("not enough money")
	ErrBalanceTooLow  = errors.New("insufficient funds to place current wager")

	ErrWrongSessionToken                    = errors.New("wrong session token")
	ErrSessionTokenExpired                  = errors.New("session token expired")
	ErrWrongFreeSpinID                      = errors.New("wrong free spin id")
	ErrBadDataGiven                         = errors.New("bad data given")
	ErrHistoryRecordNotFound                = errors.New("history record not found")
	ErrSpinGenerationCanNotBeContinued      = errors.New("spin generation can not be continued")
	ErrGameNotSupportsFreeSpins             = errors.New("game not supports free spins")
	ErrGambleAnyWinWasDisabledOnServerLevel = errors.New("gamble any win was disabled on server level")
	ErrLimitForGambleSetToZero              = errors.New("limit for gamble is set to 0")
	ErrCanNotGamble                         = errors.New("can not gamble")
	ErrUserHasDifferentCurrency             = errors.New("user_has_different_currency")

	ErrUserIsBlocked             = errors.New("user is blocked")
	ErrIntegratorCriticalFailure = errors.New("integrator critical failure")

	ErrInternalBadData = errors.New("internal bad data")

	ErrInitImpossibleBootConfig = errors.New("impossible boot config")
)

func OneOfListError[T comparable](field string, list []T) string {
	return fmt.Sprintf("%s must be one of %v", field, list)
}
