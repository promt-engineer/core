package main

import (
	"bitbucket.org/play-workspace/base-slot-server/buildvar"
	"bitbucket.org/play-workspace/base-slot-server/internal/roulette"
	"bitbucket.org/play-workspace/base-slot-server/pkg/app"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/constants"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"bitbucket.org/play-workspace/base-slot-server/tests/stress"
	"go.uber.org/zap"
	"net/url"
)

func main() {
	application, err := app.NewApp("config.yml", buildvar.Tag, buildvar.Debug, buildvar.IsCheatsAvailable, roulette.GameBoot)
	if err != nil {
		panic(err)
	}

	logger := application.Ctn().Get(constants.LoggerName).(*zap.Logger)
	url, err := url.Parse("https://dev.bf.heronbyte.com/smashing-hot-94/api")
	if err != nil {
		panic(err)
	}

	spammer := stress.NewSpammer(helper{}, engine.GetFromContainer().SpinFactory, logger, *url, "smashing-hot-94")

	rep, err := spammer.Run(2000, 3)
	if err != nil {
		logger.Sugar().Info(err)
	}

	if len(rep.Errs) > 0 {
		logger.Sugar().Info(len(rep.Errs))
		logger.Sugar().Info(rep.Errs)
	} else {
		logger.Sugar().Infof("init lat: %v, wager lat: %v, update lat: %v",
			rep.AvgInitLatency(), rep.AvgWagerLatency(), rep.AvgUpdateLatency())
	}
}

type helper struct {
}

func (h helper) WagerParams() interface{} {
	return nil
}

func (h helper) UpdateParams(spin engine.Spin, restoring engine.RestoringIndexes) interface{} {
	return map[string]interface{}{"base_spin_index": 1}
}
