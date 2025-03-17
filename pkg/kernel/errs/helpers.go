package errs

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/cryptolut_rgs"
	"bitbucket.org/play-workspace/base-slot-server/pkg/history"
	"bitbucket.org/play-workspace/base-slot-server/pkg/overlord"
)

var translateOverlordMap = map[error]error{
	overlord.ErrMarshaling:                ErrBadDataGiven,
	overlord.ErrBalanceTooLow:             ErrNotEnoughMoney,
	overlord.ErrWrongSessionToken:         ErrWrongSessionToken,
	overlord.ErrWrongFreeSpinID:           ErrWrongFreeSpinID,
	overlord.ErrSessionTokenExpired:       ErrSessionTokenExpired,
	overlord.ErrUserHasDifferentCurrency:  ErrUserHasDifferentCurrency,
	overlord.ErrUserIsBlocked:             ErrUserIsBlocked,
	overlord.ErrIntegratorCriticalFailure: ErrIntegratorCriticalFailure,

	cryptolut_rgs.ErrCLSessionNotFound:     ErrCLSessionNotFound,
	cryptolut_rgs.ErrCLNotConfigured:       ErrCLNotConfigured,
	cryptolut_rgs.ErrCLAccessDenied:        ErrCLAccessDenied,
	cryptolut_rgs.ErrCLValidationError:     ErrCLValidationError,
	cryptolut_rgs.ErrCLInternalServerError: ErrCLInternalServerError,
	cryptolut_rgs.ErrCLCurrencyNotFound:    ErrCLCurrencyNotFound,
	cryptolut_rgs.ErrCLCountryNotFound:     ErrCLCountryNotFound,
	cryptolut_rgs.ErrCLNotAvailable:        ErrCLNotAvailable,
	cryptolut_rgs.ErrCLSessionTimeout:      ErrCLSessionTimeout,
	cryptolut_rgs.ErrCLSessionDuplicate:    ErrCLSessionDuplicate,
	cryptolut_rgs.ErrCLSessionReopened:     ErrCLSessionReopened,
}

var translateHistoryMap = map[error]error{
	history.ErrSpinNotFound: ErrHistoryRecordNotFound,
}

func TranslateOverlordErr(err error) error {
	validationErr, ok := err.(overlord.ValidationError)
	if ok {
		return InternalValidationError{Err: validationErr}
	}

	res, ok := translateOverlordMap[err]
	if !ok {
		return err
	}

	return res
}

func TranslateHistoryErr(err error) error {
	validationErr, ok := err.(overlord.ValidationError)
	if ok {
		return InternalValidationError{Err: validationErr}
	}

	res, ok := translateHistoryMap[err]
	if !ok {
		return err
	}

	return res
}
