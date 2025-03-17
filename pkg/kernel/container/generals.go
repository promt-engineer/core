package container

import (
	"fmt"
	"time"

	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/config"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/constants"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/facade"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/services"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/validator"
	"bitbucket.org/play-workspace/gocommon/tracer"
	"github.com/go-co-op/gocron"
	"github.com/sarulabs/di"
	"go.uber.org/zap"
)

func BuildGenerals(configPath string, isDebug bool) []di.Def {

	return append(buildGenerals(configPath),
		[]di.Def{
			{
				Name: constants.LoggerName,
				Build: func(ctn di.Container) (interface{}, error) {
					logger, err := zap.NewProduction()
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					if isDebug || cfg.EngineConfig.Debug {
						logger, err = zap.NewDevelopment()
					}

					if err != nil {
						return nil, fmt.Errorf("can't initialize zap logger: %v", err)
					}

					zap.ReplaceGlobals(logger)

					return logger, nil
				},
			},
		}...,
	)
}

func NewBuildGenerals(configPath string) []di.Def {
	return append(buildGenerals(configPath),
		[]di.Def{
			{
				Name: constants.LoggerName,
				Build: func(ctn di.Container) (interface{}, error) {
					logger, err := zap.NewProduction()
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					if cfg.EngineConfig.Debug {
						logger, err = zap.NewDevelopment()
					}

					if err != nil {
						return nil, fmt.Errorf("can't initialize zap logger: %v", err)
					}

					zap.ReplaceGlobals(logger)

					return logger, nil
				},
			},
		}...,
	)
}

func buildGenerals(configPath string) []di.Def {
	return []di.Def{

		{
			Name: constants.ConfigName,
			Build: func(ctn di.Container) (interface{}, error) {
				return config.New(configPath)
			},
		},
		{
			Name: constants.ValidatorName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				return validator.New(cfg.ConstantsConfig)
			},
		},
		{
			Name: constants.SchedulerName,
			Build: func(ctn di.Container) (interface{}, error) {
				s := gocron.NewScheduler(time.UTC)
				s.StartAsync()

				return s, nil
			},
			Close: func(obj interface{}) error {
				obj.(*gocron.Scheduler).Clear()
				zap.S().Info("Scheduler stopped")

				return nil
			},
		},
		{
			Name: constants.TracerName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				return tracer.NewTracer(cfg.TracerConfig)
			},
		},
		{
			Name: constants.FacadeName,
			Build: func(ctn di.Container) (interface{}, error) {
				validationEngine := ctn.Get(constants.ValidatorName).(*validator.Validator)
				gameFlow := ctn.Get(constants.GameFlowServiceName).(*services.GameFlowService)
				history := ctn.Get(constants.HistoryServiceName).(*services.HistoryService)
				freeSpin := ctn.Get(constants.FreeSpinServiceName).(*services.FreeSpinService)
				cheats := ctn.Get(constants.CheatsServiceName).(*services.CheatsService)

				return facade.NewFacade(validationEngine, gameFlow, history, freeSpin, cheats), nil
			},
		},
	}
}
