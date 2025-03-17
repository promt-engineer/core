package websocket

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/validator"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Response struct {
	Status  int                    `json:"status"`
	Success bool                   `json:"success"`
	Meta    map[string]interface{} `json:"meta"`
	Data    interface{}            `json:"data"`
}

func new(status int, meta map[string]interface{}, data interface{}) *Response {
	success := false
	if status == StatusSuccess {
		success = true
	}

	response := &Response{
		Status:  status,
		Success: success,
		Meta:    meta,
		Data:    data,
	}

	if v, ok := data.(error); ok {
		response.Data = v.Error()
	}

	return response
}

func OK(data interface{}, uuid uuid.UUID) *Response {
	meta := map[string]interface{}{"uuid": uuid}

	return new(StatusSuccess, meta, data)
}

func OKNoContent(uuid uuid.UUID) *Response {
	meta := map[string]interface{}{"uuid": uuid}

	return new(StatusSuccess, meta, "no content")
}

func BadRequest(data interface{}) *Response {
	return new(StatusBadRequest, nil, data)
}

func NotFound(uuid uuid.UUID) *Response {
	meta := map[string]interface{}{"uuid": uuid}

	return new(StatusActionNotFound, meta, nil)
}

func Unauthorized(data interface{}, uuid uuid.UUID) *Response {
	meta := map[string]interface{}{"uuid": uuid}

	return new(StatusUnauthorized, meta, data)
}

func PaymentRequired(data interface{}, uuid uuid.UUID) *Response {
	meta := map[string]interface{}{"uuid": uuid}

	return new(StatusPaymentRequirement, meta, data)
}

func Forbidden(data interface{}, uuid uuid.UUID) *Response {
	meta := map[string]interface{}{"uuid": uuid}

	return new(StatusForbidden, meta, data)
}

func Conflict(data interface{}, uuid uuid.UUID) *Response {
	meta := map[string]interface{}{"uuid": uuid}

	return new(StatusConflict, meta, data)
}

func SessionExpired(data interface{}, uuid uuid.UUID) *Response {
	meta := map[string]interface{}{"uuid": uuid}

	return new(StatusSessionExpired, meta, data)
}

func ServerError(data interface{}, uuid uuid.UUID) *Response {
	meta := map[string]interface{}{"uuid": uuid}

	zap.S().Error(data)

	return new(StatusInternalError, meta, data)
}

func ValidationFailed(err error) *Response {
	data := []string{}

	for _, taggedError := range validator.CheckValidationErrors(err) {
		e := taggedError.Err
		data = append(data, e.Error())
	}

	return new(StatusValidationFailed, nil, data)
}
