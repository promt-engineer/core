package handlers

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/entities"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/facade"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type gameFlowHandler struct {
	facade *facade.Facade
}

func NewGameFlowHandler(facade *facade.Facade) http.Handler {
	return &gameFlowHandler{facade: facade}
}

func (h *gameFlowHandler) Shutdown() {}

func (h *gameFlowHandler) Register(router *gin.RouterGroup) {
	core := router.Group("core")

	core.POST("state", h.initState)
	core.POST("wager", h.wager)
	core.POST("gamble_any_win", h.gambleAnyWin)
	core.POST("keep_generating", h.keepGenerating)
	core.GET("spins_history", h.history)
	core.POST("spin_indexes/update", h.updateSpinIndexes)

	core.GET("free_spins/get", h.getFreeSpins)
	core.GET("free_spins/cancel", h.cancelFreeSpins)

	core.GET("free_spins/get_with_integrator_bet", h.getFreeSpinsWithIntegratorBet)
	core.GET("free_spins/cancel_with_integrator_bet", h.cancelFreeSpinsWithIntegratorBet)
}

type jsonInterface interface{}

func (h *gameFlowHandler) initState(ctx *gin.Context) {
	payload, err := bindBody(ctx)
	if err != nil {
		zap.S().Error("initState: ", err)
		http.BadRequest(ctx, err, nil)

		return
	}

	gameState, err := h.facade.InitState(ctx.Request.Context(), payload)
	if err != nil {
		handleServiceError(ctx, err)

		return
	}

	http.OK(ctx, gameState, nil)
}

func (h *gameFlowHandler) wager(ctx *gin.Context) {
	payload, err := bindBody(ctx)
	if err != nil {
		zap.S().Error("wager: ", err)
		http.BadRequest(ctx, err, nil)

		return
	}

	md, err := getMetaData(ctx, payload)
	if err != nil {
		zap.S().Error("getMetaData: ", err)
		handleServiceError(ctx, err)

		return
	}

	gameState, err := h.facade.Wager(ctx.Request.Context(), payload, md)
	if err != nil {
		handleServiceError(ctx, err)

		return
	}

	http.OK(ctx, gameState, nil)
}

func (h *gameFlowHandler) gambleAnyWin(ctx *gin.Context) {
	payload, err := bindBody(ctx)
	if err != nil {
		zap.S().Error("gamble any win: ", err)
		http.BadRequest(ctx, err, nil)

		return
	}

	md, err := getMetaData(ctx, payload)
	if err != nil {
		zap.S().Error("getMetaData: ", err)
		handleServiceError(ctx, err)

		return
	}

	gameState, err := h.facade.GambleAnyWin(ctx.Request.Context(), payload, md)
	if err != nil {
		handleServiceError(ctx, err)

		return
	}

	http.OK(ctx, gameState, nil)
}

func (h *gameFlowHandler) keepGenerating(ctx *gin.Context) {
	payload, err := bindBody(ctx)
	if err != nil {
		zap.S().Error("keep generating: ", err)
		http.BadRequest(ctx, err, nil)

		return
	}

	md, err := getMetaData(ctx, payload)
	if err != nil {
		zap.S().Error("getMetaData: ", err)
		handleServiceError(ctx, err)

		return
	}

	gameState, err := h.facade.KeepGenerating(ctx.Request.Context(), payload, md)
	if err != nil {
		handleServiceError(ctx, err)

		return
	}

	http.OK(ctx, gameState, nil)
}

func (h *gameFlowHandler) history(ctx *gin.Context) {
	payload := queryToMap(ctx)

	pagination, err := h.facade.Paginate(ctx.Request.Context(), payload)
	if err != nil {
		handleServiceError(ctx, err)

		return
	}

	http.OK(ctx, pagination, nil)
}

func (h *gameFlowHandler) updateSpinIndexes(ctx *gin.Context) {
	payload, err := bindBody(ctx)
	if err != nil {
		zap.S().Error("upd spin index: ", err)
		http.BadRequest(ctx, err, nil)

		return
	}

	md, err := getMetaData(ctx, payload)
	if err != nil {
		zap.S().Error("getMetaData: ", err)
		handleServiceError(ctx, err)

		return
	}

	if err := h.facade.UpdateSpinIndexes(ctx.Request.Context(), payload, md); err != nil {
		handleServiceError(ctx, err)

		return
	}

	http.OK(ctx, nil, nil)
}

func (h *gameFlowHandler) getFreeSpins(ctx *gin.Context) {
	payload := queryToMap(ctx)

	freeSpins, err := h.facade.FreeSpins(ctx.Request.Context(), payload)
	if err != nil {
		zap.S().Error("get fs: ", err)
		handleServiceError(ctx, err)

		return
	}

	http.OK(ctx, freeSpins, nil)
}

func (h *gameFlowHandler) cancelFreeSpins(ctx *gin.Context) {
	payload := queryToMap(ctx)

	if err := h.facade.CancelSpins(ctx.Request.Context(), payload); err != nil {
		zap.S().Error("cancel fs: ", err)
		handleServiceError(ctx, err)

		return
	}

	http.OKNoContent(ctx)
}

func (h *gameFlowHandler) getFreeSpinsWithIntegratorBet(ctx *gin.Context) {
	payload := queryToMap(ctx)

	freeSpins, err := h.facade.FreeSpinsWithIntegratorBet(ctx.Request.Context(), payload)
	if err != nil {
		zap.S().Error("get fs: ", err)
		handleServiceError(ctx, err)

		return
	}

	http.OK(ctx, freeSpins, nil)
}

func (h *gameFlowHandler) cancelFreeSpinsWithIntegratorBet(ctx *gin.Context) {
	payload := queryToMap(ctx)

	if err := h.facade.CancelSpinsWithIntegratorBet(ctx.Request.Context(), payload); err != nil {
		zap.S().Error("cancel fs: ", err)
		handleServiceError(ctx, err)

		return
	}

	http.OKNoContent(ctx)
}

func getMetaData(ctx *gin.Context, parsedRequest interface{}) (*entities.PlayerMetaData, error) {
	requestBody, err := json.Marshal(parsedRequest)
	if err != nil {
		return nil, err
	}

	return entities.NewPlayerMetaDataFromCtx(ctx, requestBody), nil
}
