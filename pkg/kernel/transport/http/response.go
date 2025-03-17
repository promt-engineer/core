package http

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	CodeSessionExpired = 419
)

type Response struct {
	Status  int                    `json:"status"`
	Success bool                   `json:"success"`
	Meta    map[string]interface{} `json:"meta"`
	Data    interface{}            `json:"data"`
}

type GameHubErrorResponse struct {
	Status int               `json:"status"`
	Error  errs.GameHubError `json:"error"`
}

func new(status int, meta map[string]interface{}, data interface{}) *Response {
	success := false
	if status >= 200 && status <= 299 {
		success = true
	}

	response := &Response{
		Status:  status,
		Success: success,
		Meta:    meta,
		Data:    data,
	}

	if response.Data == nil {
		response.Data = http.StatusText(status)
	}

	if v, ok := data.(error); ok {
		response.Data = v.Error()
	}

	return response
}

func sendError(ctx *gin.Context, status int, err error, meta map[string]interface{}) {
	zap.S().Error(err)

	if gameHubErr, ok := errs.MapErrorToGameHub(err); ok {
		ctx.JSON(status, GameHubErrorResponse{
			Status: status,
			Error:  gameHubErr,
		})
		return
	}

	r := new(status, meta, err)
	ctx.JSON(status, r)
}

func OK(ctx *gin.Context, data interface{}, meta map[string]interface{}) {
	r := new(http.StatusOK, meta, data)
	ctx.JSON(r.Status, r)
}

func OKNoContent(ctx *gin.Context) {
	ctx.Data(http.StatusNoContent, gin.MIMEHTML, nil)
}

func BadRequest(ctx *gin.Context, err error, meta map[string]interface{}) {
	sendError(ctx, http.StatusBadRequest, err, meta)
}

func Unauthorized(ctx *gin.Context, err error, meta map[string]interface{}) {
	sendError(ctx, http.StatusUnauthorized, err, meta)
}

func PaymentRequired(ctx *gin.Context, err error, meta map[string]interface{}) {
	sendError(ctx, http.StatusPaymentRequired, err, meta)
}

func Forbidden(ctx *gin.Context, err error, meta map[string]interface{}) {
	sendError(ctx, http.StatusForbidden, err, meta)
}

func NotFound(ctx *gin.Context, err error, meta map[string]interface{}) {
	sendError(ctx, http.StatusNotFound, err, meta)
}

func ServerError(ctx *gin.Context, err error, meta map[string]interface{}) {
	sendError(ctx, http.StatusInternalServerError, err, meta)
}

func ServiceUnavailableError(ctx *gin.Context, err error, meta map[string]interface{}) {
	sendError(ctx, http.StatusServiceUnavailable, err, meta)
}

func ValidationFailed(ctx *gin.Context, err error) {
	sendError(ctx, http.StatusUnprocessableEntity, err, nil)
}

func Conflict(ctx *gin.Context, err error, meta map[string]interface{}) {
	sendError(ctx, http.StatusConflict, err, meta)
}

func CustomCode(code int) func(ctx *gin.Context, err error, meta map[string]interface{}) {
	return func(ctx *gin.Context, err error, meta map[string]interface{}) {
		sendError(ctx, code, err, meta)
	}
}

func SessionExpired(ctx *gin.Context, err error, meta map[string]interface{}) {
	sendError(ctx, CodeSessionExpired, err, meta)
}
