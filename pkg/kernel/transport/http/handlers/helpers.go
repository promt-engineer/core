package handlers

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/errs"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/http"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
)

var errorMap = map[error]func(ctx *gin.Context, err error, meta map[string]interface{}){
	errs.ErrSessionTokenExpired: http.SessionExpired,
	errs.ErrWrongSessionToken:   http.Unauthorized,

	errs.ErrHistoryRecordNotFound: http.Conflict,
	errs.ErrLastSpinWasNotShown:   http.Conflict,
	errs.ErrWrongFreeSpinID:       http.Conflict,
	errs.ErrNotEnoughMoney:        http.PaymentRequired,
	errs.ErrBalanceTooLow:         http.PaymentRequired,

	errs.ErrUserIsBlocked:             http.Forbidden,
	errs.ErrUserHasDifferentCurrency:  http.Conflict,
	errs.ErrIntegratorCriticalFailure: http.ServiceUnavailableError,

	errs.ErrCLSessionNotFound:     http.CustomCode(451),
	errs.ErrCLNotConfigured:       http.CustomCode(452),
	errs.ErrCLAccessDenied:        http.CustomCode(453),
	errs.ErrCLValidationError:     http.CustomCode(454),
	errs.ErrCLInternalServerError: http.ServerError,
	errs.ErrCLCurrencyNotFound:    http.CustomCode(455),
	errs.ErrCLCountryNotFound:     http.CustomCode(456),
	errs.ErrCLNotAvailable:        http.CustomCode(457),
	errs.ErrCLSessionTimeout:      http.CustomCode(458),
	errs.ErrCLSessionDuplicate:    http.CustomCode(459),
	errs.ErrCLSessionReopened:     http.CustomCode(460),
}

func handleServiceError(ctx *gin.Context, err error) {
	internalValidationError, ok := err.(errs.InternalValidationError)
	if ok {
		http.ValidationFailed(ctx, internalValidationError.Err)

		return
	}

	fn, ok := errorMap[err]
	if !ok {
		http.ServerError(ctx, err, nil)

		return
	}

	fn(ctx, err, nil)
}

func queryToMap(ctx *gin.Context) map[string]interface{} {
	paramMap := make(map[string]interface{}, 0)

	for k, v := range ctx.Request.URL.Query() {
		if len(v) == 1 && len(v[0]) != 0 {
			paramMap[k] = v[0]

			i, err := strconv.Atoi(v[0])
			if err == nil {
				paramMap[k] = i
			}
		} else {
			continue
		}
	}

	return paramMap
}

func bindBody(ctx *gin.Context) (jsonInterface, error) {
	payload := new(jsonInterface)
	if err := ctx.ShouldBind(payload); err != nil {
		return nil, err
	}

	if reflect.ValueOf(payload).Elem().IsNil() {
		return nil, errs.ErrBadDataGiven
	}

	return *payload, nil
}
