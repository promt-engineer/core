package app

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/config"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/constants"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/container"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	utils2 "bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine/utils"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine/utils/volatility"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/errs"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/services"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/http"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/validator"
	"bitbucket.org/play-workspace/base-slot-server/pkg/rng"
	"bitbucket.org/play-workspace/base-slot-server/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/sarulabs/di"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

type App struct {
	ctn di.Container
	ctx context.Context
	wg  *sync.WaitGroup
}

type GameBootstrapV2 func(rand rng.Client, volatility volatility.Type, rtp float64) *engine.Bootstrap

func New(configPath string, fn GameBootstrapV2) (*App, error) {
	app := &App{
		ctx: context.Background(),
		wg:  &sync.WaitGroup{},
	}

	app.ctn = container.NewBuild(app.ctx, app.wg, configPath)

	logger := app.ctn.Get(constants.LoggerName).(*zap.Logger)
	logger.Info("Building application...")

	cfg := app.ctn.Get(constants.ConfigName).(*config.Config)

	rand := utils2.GetRNG(app.ctn, cfg.EngineConfig)

	rtp, err := strconv.ParseFloat(cfg.EngineConfig.RTP, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing rtp: %w", err)
	}

	vol, err := volatility.VolFromStr(cfg.EngineConfig.Volatility)
	if err != nil {
		return nil, fmt.Errorf("error parsing volatility: %w", err)
	}

	var boot = fn(rand, vol, rtp)

	if boot.HistoryHandlingType != engine.SequentialRestoring && boot.GambleAnyWinFeature {
		return nil, errs.ErrInitImpossibleBootConfig
	}

	engine.PutInContainer(boot)

	return app, nil
}

func (app *App) Ctn() di.Container {
	return app.ctn
}

func (app *App) Run() error {
	now := time.Now()

	server := app.ctn.Get(constants.ServerName).(*http.Server)
	binding.Validator = app.ctn.Get(constants.ValidatorName).(*validator.Validator)

	go server.Run()

	zap.S().Infof("Up and running (%s)", time.Since(now))
	zap.S().Infof("Got %s signal. Shutting down...", <-utils.WaitTermSignal())

	if err := server.Shutdown(context.Background()); err != nil {
		zap.S().Errorf("Error stopping server: %s", err)
	}

	app.wg.Wait()
	zap.S().Info("Service stopped.")

	return nil
}

func (app *App) RunOrSimulate(keepGenerate services.KeepGenerateWrapper) error {
	cfg := app.ctn.Get(constants.ConfigName).(*config.Config)

	// Run the application if simulator config is not available
	if cfg.SimulatorConfig == nil {
		return app.Run()
	}

	simService := app.ctn.Get(constants.SimulatorServiceName).(*services.SimulatorService)
	if simService == nil {
		return fmt.Errorf("simulator service is not available")
	}

	if keepGenerate != nil {
		simService.WithKeepGenerate(keepGenerate)
	}

	return simService.SimulateV2(cfg.SimulatorConfig, cfg.EngineConfig.RTP, cfg.EngineConfig.Volatility)
}
