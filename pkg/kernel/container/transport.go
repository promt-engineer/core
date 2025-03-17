package container

import (
	"context"
	"sync"

	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/config"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/constants"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/facade"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/http"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/websocket"
	"bitbucket.org/play-workspace/gocommon/tracer"
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"
)

func BuildTransport(ctx context.Context, wg *sync.WaitGroup) []di.Def {
	return []di.Def{
		{
			Name: constants.WebsocketServerName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)
				fcd := ctn.Get(constants.FacadeName).(*facade.Facade)
				gameFlowHandler := ctn.Get(constants.WSGameFlowHandlerName).(websocket.Handler)
				cheatsHandler := ctn.Get(constants.WSCheatsHandlerName).(websocket.Handler)
				tr := ctn.Get(constants.TracerName).(*tracer.JaegerTracer)

				return websocket.NewServer(cfg.WebsocketConfig, fcd, tr, gameFlowHandler, cheatsHandler), nil
			},
		},
		{
			Name: constants.ServerName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				publicHandlers := []http.Handler{
					ctn.Get(constants.HTTPMetaHandlerName).(http.Handler),
					ctn.Get(constants.HTTPMetricsHandlerName).(http.Handler),
					ctn.Get(constants.HTTPCheatsHandlerName).(http.Handler),
					ctn.Get(constants.HTTPSimulatorHandlerName).(http.Handler),
				}

				boot := engine.GetFromContainer()

				if boot.WebsocketTransport {
					publicHandlers = append(publicHandlers, ctn.Get(constants.HTTPWSHandlerName).(http.Handler))
				}

				if boot.HTTPTransport {
					publicHandlers = append(publicHandlers, ctn.Get(constants.HTTPGameFlowHandlerName).(http.Handler))
				}

				var privateHandlers []http.Handler

				var middlewares = []func(ctx *gin.Context){
					ctn.Get(constants.HTTPCorsMiddlewareName).(func(ctx *gin.Context)),
					ctn.Get(constants.HTTPSessionMiddlewareName).(func(ctx *gin.Context)),
					ctn.Get(constants.HTTPSessionMuMiddlewareName).(func(ctx *gin.Context)),
					ctn.Get(constants.HTTPTraceMiddlewareName).(func(ctx *gin.Context)),
				}

				return http.New(ctx, wg, cfg.ServerConfig, cfg.ConstantsConfig, publicHandlers, privateHandlers, middlewares), nil
			},
		},
	}
}
