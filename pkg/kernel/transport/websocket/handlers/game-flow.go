package handlers

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/facade"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/websocket"
)

type gameFlowHandler struct {
	facade *facade.Facade
}

func NewGameFlowHandler(facade *facade.Facade) websocket.Handler {
	return &gameFlowHandler{facade: facade}
}

func (h *gameFlowHandler) Shutdown() {}

func (h *gameFlowHandler) Register(r *websocket.Router) {
	r.Accept(ActionState, h.state)
	r.Accept(ActionWager, h.wager)
	r.Accept(ActionGambleAnyWin, h.gambleAnyWin)
	r.Accept(ActionKeepGenerating, h.keepGenerating)
	r.Accept(ActionSpinsHistory, h.spinsHistory)
	r.Accept(ActionUpdateSpinIndexes, h.updateSpinIndexes)
	r.Accept(ActionGetFreeSpins, h.getFreeSpins)
	r.Accept(ActionCancelFreeSpins, h.cancelFreeSpins)
	r.Accept(ActionGetFreeSpinsWithIntegratorBet, h.getFreeSpinsWithIntegratorBet)
	r.Accept(ActionCancelFreeSpinsWithIntegratorBet, h.cancelFreeSpinsWithIntegratorBet)
}

func (h *gameFlowHandler) state(bag websocket.HandlerBag) {
	gameState, err := h.facade.InitState(bag.Ctx, bag.Payload)
	if err != nil {
		handleServiceError(bag.ResponsePipeline, err, bag.UUID)

		return
	}

	bag.ResponsePipeline <- websocket.OK(gameState, bag.UUID)
}

func (h *gameFlowHandler) wager(bag websocket.HandlerBag) {
	gameState, err := h.facade.Wager(bag.Ctx, bag.Payload, bag.PlayerMetaData)
	if err != nil {
		handleServiceError(bag.ResponsePipeline, err, bag.UUID)

		return
	}

	bag.ResponsePipeline <- websocket.OK(gameState, bag.UUID)
}

func (h *gameFlowHandler) gambleAnyWin(bag websocket.HandlerBag) {
	gameState, err := h.facade.GambleAnyWin(bag.Ctx, bag.Payload, bag.PlayerMetaData)
	if err != nil {
		handleServiceError(bag.ResponsePipeline, err, bag.UUID)

		return
	}

	bag.ResponsePipeline <- websocket.OK(gameState, bag.UUID)
}

func (h *gameFlowHandler) keepGenerating(bag websocket.HandlerBag) {
	gameState, err := h.facade.KeepGenerating(bag.Ctx, bag.Payload, bag.PlayerMetaData)
	if err != nil {
		handleServiceError(bag.ResponsePipeline, err, bag.UUID)

		return
	}

	bag.ResponsePipeline <- websocket.OK(gameState, bag.UUID)
}

func (h *gameFlowHandler) getFreeSpins(bag websocket.HandlerBag) {
	freeSpins, err := h.facade.FreeSpins(bag.Ctx, bag.Payload)
	if err != nil {
		handleServiceError(bag.ResponsePipeline, err, bag.UUID)

		return
	}

	bag.ResponsePipeline <- websocket.OK(freeSpins, bag.UUID)
}

func (h *gameFlowHandler) cancelFreeSpins(bag websocket.HandlerBag) {
	if err := h.facade.CancelSpins(bag.Ctx, bag.Payload); err != nil {
		handleServiceError(bag.ResponsePipeline, err, bag.UUID)

		return
	}

	bag.ResponsePipeline <- websocket.OKNoContent(bag.UUID)
}

func (h *gameFlowHandler) getFreeSpinsWithIntegratorBet(bag websocket.HandlerBag) {
	freeSpins, err := h.facade.FreeSpinsWithIntegratorBet(bag.Ctx, bag.Payload)
	if err != nil {
		handleServiceError(bag.ResponsePipeline, err, bag.UUID)

		return
	}

	bag.ResponsePipeline <- websocket.OK(freeSpins, bag.UUID)
}

func (h *gameFlowHandler) cancelFreeSpinsWithIntegratorBet(bag websocket.HandlerBag) {
	if err := h.facade.CancelSpinsWithIntegratorBet(bag.Ctx, bag.Payload); err != nil {
		handleServiceError(bag.ResponsePipeline, err, bag.UUID)

		return
	}

	bag.ResponsePipeline <- websocket.OKNoContent(bag.UUID)
}

func (h *gameFlowHandler) spinsHistory(bag websocket.HandlerBag) {
	pagination, err := h.facade.Paginate(bag.Ctx, bag.Payload)
	if err != nil {
		handleServiceError(bag.ResponsePipeline, err, bag.UUID)

		return
	}

	bag.ResponsePipeline <- websocket.OK(pagination, bag.UUID)
}

func (h *gameFlowHandler) updateSpinIndexes(bag websocket.HandlerBag) {
	if err := h.facade.UpdateSpinIndexes(bag.Ctx, bag.Payload, bag.PlayerMetaData); err != nil {
		handleServiceError(bag.ResponsePipeline, err, bag.UUID)

		return
	}

	bag.ResponsePipeline <- websocket.OK(nil, bag.UUID)
}
