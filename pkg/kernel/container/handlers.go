package container

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/config"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/constants"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/facade"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/services"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/http/handlers"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/websocket"
	wsHandlers "bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/websocket/handlers"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/validator"
	"github.com/sarulabs/di"
)

func BuildHandlers(buildTag string, isCheatsAvailable bool) []di.Def {
	defs := append(buildHandlers(),
		[]di.Def{
			{
				Name: constants.HTTPMetaHandlerName,
				Build: func(ctn di.Container) (interface{}, error) {
					return handlers.NewMetaHandler(buildTag), nil
				},
			},
			{
				Name: constants.HTTPCheatsHandlerName,
				Build: func(ctn di.Container) (interface{}, error) {
					fcd := ctn.Get(constants.FacadeName).(*facade.Facade)
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					return handlers.NewCheatsHandler(fcd, isCheatsAvailable || cfg.EngineConfig.IsCheatsAvailable), nil
				},
			},
			{
				Name: constants.WSCheatsHandlerName,
				Build: func(ctn di.Container) (interface{}, error) {
					fcd := ctn.Get(constants.FacadeName).(*facade.Facade)
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					return wsHandlers.NewCheatsHandler(fcd, isCheatsAvailable || cfg.EngineConfig.IsCheatsAvailable), nil
				},
			},
		}...,
	)

	return defs
}

func NewBuildHandlers() []di.Def {
	defs := append(buildHandlers(),
		[]di.Def{
			{
				Name: constants.HTTPMetaHandlerName,
				Build: func(ctn di.Container) (interface{}, error) {
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					return handlers.NewMetaHandler(cfg.EngineConfig.BuildVersion), nil
				},
			},
			{
				Name: constants.HTTPCheatsHandlerName,
				Build: func(ctn di.Container) (interface{}, error) {
					fcd := ctn.Get(constants.FacadeName).(*facade.Facade)
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					return handlers.NewCheatsHandler(fcd, cfg.EngineConfig.IsCheatsAvailable), nil
				},
			},
			{
				Name: constants.WSCheatsHandlerName,
				Build: func(ctn di.Container) (interface{}, error) {
					fcd := ctn.Get(constants.FacadeName).(*facade.Facade)
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					return wsHandlers.NewCheatsHandler(fcd, cfg.EngineConfig.IsCheatsAvailable), nil
				},
			},
		}...,
	)

	return defs
}

func buildHandlers() []di.Def {
	return []di.Def{
		{
			Name: constants.HTTPGameFlowHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				fcd := ctn.Get(constants.FacadeName).(*facade.Facade)

				return handlers.NewGameFlowHandler(fcd), nil
			},
		},

		{
			Name: constants.HTTPMetricsHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				return handlers.NewMetricsHandler(), nil
			},
		},
		{
			Name: constants.HTTPWSHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				websocketServer := ctn.Get(constants.WebsocketServerName).(*websocket.Server)

				return handlers.NewWebsocketHandler(websocketServer), nil
			},
		},

		{
			Name: constants.WSGameFlowHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				fcd := ctn.Get(constants.FacadeName).(*facade.Facade)

				return wsHandlers.NewGameFlowHandler(fcd), nil
			},
		},
		{
			Name: constants.HTTPSimulatorHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				validatorEngine := ctn.Get(constants.ValidatorName).(*validator.Validator)
				simulatorService := ctn.Get(constants.SimulatorServiceName).(*services.SimulatorService)
				return handlers.NewSimulatorHandler(ctn, validatorEngine, simulatorService), nil
			},
		},
	}
}
