package handlers

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/http"
	"github.com/gin-gonic/gin"
)

type metaHandler struct {
	Tag string
}

type HealthResponse struct {
	Status string
}

type InfoResponse struct {
	Tag string
	IP  string
}

func NewMetaHandler(buildTag string) http.Handler {
	return &metaHandler{Tag: buildTag}
}

func (h *metaHandler) Shutdown() {
}

func (h *metaHandler) Register(route *gin.RouterGroup) {
	route.GET("health", h.health)
	route.Any("info", h.info)
}

// @Summary Check health.
// @Tags meta
// @Consume application/json
// @Description Check service health.
// @Accept  json
// @Produce  json
// @Success 200  {object} responses.HealthResponse
// @Router /health [get].
func (h *metaHandler) health(ctx *gin.Context) {
	http.OK(ctx, HealthResponse{Status: "success"}, nil)
}

// @Summary Check tag.
// @Tags meta
// @Consume application/json
// @Description Check service tag.
// @Accept  json
// @Produce  json
// @Success 200  {object} responses.InfoResponse
// @Router /info [get].
func (h *metaHandler) info(ctx *gin.Context) {
	http.OK(ctx, InfoResponse{Tag: h.Tag, IP: ctx.ClientIP()}, nil)
}
