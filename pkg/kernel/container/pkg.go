package container

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/cryptolut_rgs"
	"bitbucket.org/play-workspace/base-slot-server/pkg/history"
	"bitbucket.org/play-workspace/base-slot-server/pkg/ip2country"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/config"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/constants"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/validator"
	"bitbucket.org/play-workspace/base-slot-server/pkg/overlord"
	"bitbucket.org/play-workspace/base-slot-server/pkg/rng"
	"fmt"
	"github.com/sarulabs/di"
	"strings"
	"time"
)

func BuildPkg() []di.Def {
	return []di.Def{
		{
			Name: constants.OverlordName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				if cfg.OverlordConfig != nil {
					return overlord.NewClient(cfg.OverlordConfig)
				}

				return cryptolut_rgs.NewClient(cfg.CryptolutConfig), nil
			},
		},
		{
			Name: ip2country.IP2CountryName,
			Build: func(ctn di.Container) (interface{}, error) {
				c := ip2country.NewClientWithCache(time.Hour*2, strings.ToLower)

				go c.Start()

				return c, nil
			},
			Close: func(obj interface{}) error {
				c, ok := obj.(*ip2country.ClientWithCache)

				if !ok {
					return fmt.Errorf("can not convert %T to *ip2country.ClientWithCache", obj)
				}

				c.Stop()

				return nil
			},
		},
		{
			Name: constants.HistoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)
				vld := ctn.Get(constants.ValidatorName).(*validator.Validator)
				ip2C := ctn.Get(ip2country.IP2CountryName).(*ip2country.ClientWithCache)

				if cfg.HistoryConfig == nil && cfg.HistoryMongoDBConfig != nil {
					return history.NewMongoDBClient(cfg.HistoryMongoDBConfig, vld, ip2C)
				}

				return history.NewClient(cfg.HistoryConfig)
			},
		},
		{
			Name: constants.RNGName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				return rng.NewSimpleClient(cfg.RNGConfig)
			},
		},
		{
			Name: constants.RNGMockName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				return rng.NewMockClient(cfg.RNGConfig)
			},
		},
	}
}
