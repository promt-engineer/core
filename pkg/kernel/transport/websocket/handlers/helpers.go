package handlers

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/errs"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/websocket"
	"github.com/google/uuid"
)

var errorMap = map[error]func(data interface{}, uuid uuid.UUID) *websocket.Response{
	errs.ErrSessionTokenExpired: websocket.SessionExpired,
	errs.ErrWrongSessionToken:   websocket.Unauthorized,

	errs.ErrHistoryRecordNotFound: websocket.Conflict,
	errs.ErrWrongFreeSpinID:       websocket.Conflict,
	errs.ErrLastSpinWasNotShown:   websocket.Conflict,
	errs.ErrNotEnoughMoney:        websocket.PaymentRequired,
}

func handleServiceError(broadcaster chan *websocket.Response, err error, requestUUID uuid.UUID) {
	internalValidationError, ok := err.(errs.InternalValidationError)
	if ok {
		broadcaster <- websocket.ValidationFailed(internalValidationError.Err)

		return
	}

	fn, ok := errorMap[err]
	if !ok {
		broadcaster <- websocket.ServerError(err, requestUUID)

		return
	}

	broadcaster <- fn(err, requestUUID)
}
