package handlers

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/facade"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/websocket"
	"context"
)

type cheatsHandler struct {
	facade            *facade.Facade
	isCheatsAvailable bool
}

func NewCheatsHandler(facade *facade.Facade, isCheatsAvailable bool) websocket.Handler {
	return &cheatsHandler{
		facade:            facade,
		isCheatsAvailable: isCheatsAvailable,
	}
}

func (h *cheatsHandler) Register(r *websocket.Router) {
	if h.isCheatsAvailable {
		r.Accept(ActionAddCheats, h.cheats)
	}
}

func (h *cheatsHandler) Shutdown() {}

func (h *cheatsHandler) cheats(bag websocket.HandlerBag) {
	if err := h.facade.AddCheat(context.Background(), bag.Payload); err != nil {
		handleServiceError(bag.ResponsePipeline, err, bag.UUID)

		return
	}

	bag.ResponsePipeline <- websocket.OKNoContent(bag.UUID)
}
