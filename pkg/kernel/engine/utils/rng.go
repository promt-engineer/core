package utils

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/constants"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"bitbucket.org/play-workspace/base-slot-server/pkg/rng"
	"github.com/sarulabs/di"
)

func GetRNG(ctn di.Container, cfg *engine.Config) rng.Client {
	if cfg.MockRNG {
		return ctn.Get(constants.RNGMockName).(rng.Client)
	}

	return ctn.Get(constants.RNGName).(rng.Client)
}
