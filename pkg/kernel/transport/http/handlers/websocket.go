package handlers

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/http"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/websocket"
	"github.com/gin-gonic/gin"
)

type websocketHandler struct {
	websocketServer *websocket.Server
}

func (w *websocketHandler) Shutdown() {
	w.websocketServer.Shutdown()
}

func (w *websocketHandler) Register(router *gin.RouterGroup) {
	router.GET("ws", func(ctx *gin.Context) {
		err := w.websocketServer.ServeWS(ctx)
		if err != nil {
			handleServiceError(ctx, err)
		}
	})
}

func NewWebsocketHandler(websocketServer *websocket.Server) http.Handler {
	return &websocketHandler{websocketServer: websocketServer}
}
