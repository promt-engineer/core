package handlers

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/facade"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type cheatsHandler struct {
	facade            *facade.Facade
	isCheatsAvailable bool
}

func NewCheatsHandler(facade *facade.Facade, isCheatsAvailable bool) http.Handler {
	return &cheatsHandler{facade: facade, isCheatsAvailable: isCheatsAvailable}
}

func (h *cheatsHandler) Register(router *gin.RouterGroup) {
	if h.isCheatsAvailable {
		router.POST("cheats", h.cheats)
	}
}

func (h *cheatsHandler) Shutdown() {}

func (h *cheatsHandler) cheats(ctx *gin.Context) {
	payload, err := bindBody(ctx)
	if err != nil {
		zap.S().Error("cheats", err)
		http.BadRequest(ctx, err, nil)

		return
	}

	if err := h.facade.AddCheat(ctx.Request.Context(), payload); err != nil {
		handleServiceError(ctx, err)

		return
	}

	http.OKNoContent(ctx)
}
