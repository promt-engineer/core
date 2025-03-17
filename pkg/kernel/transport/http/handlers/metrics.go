package handlers

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metricsHandler struct{}

func NewMetricsHandler() http.Handler {
	return &metricsHandler{}
}

func (m metricsHandler) Shutdown() {}

func (m metricsHandler) Register(router *gin.RouterGroup) {
	metrics := router.Group("metrics")

	fn := promhttp.Handler()

	metrics.GET("/", func(ginCtx *gin.Context) {
		fn.ServeHTTP(ginCtx.Writer, ginCtx.Request)
	})
}
