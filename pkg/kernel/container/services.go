package container

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/history"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/constants"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/services"
	"bitbucket.org/play-workspace/base-slot-server/pkg/overlord"
	"github.com/sarulabs/di"
)

func BuildServices() []di.Def {
	return []di.Def{
		{
			Name: constants.GameFlowServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				lordClint := ctn.Get(constants.OverlordName).(overlord.Client)
				historySrv := ctn.Get(constants.HistoryServiceName).(*services.HistoryService)
				cheatsSrv := ctn.Get(constants.CheatsServiceName).(*services.CheatsService)

				return services.NewGameFlowService(lordClint, historySrv, cheatsSrv), nil
			},
		},
		{
			Name: constants.HistoryServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				historyClint := ctn.Get(constants.HistoryName).(history.Client)

				return services.NewHistoryService(historyClint), nil
			},
		},
		{
			Name: constants.SimulatorServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				return services.NewSimulatorService(), nil
			},
		},
		{
			Name: constants.FreeSpinServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				lordClint := ctn.Get(constants.OverlordName).(overlord.Client)

				return services.NewFreeSpinService(lordClint), nil
			},
		},
		{
			Name: constants.CheatsServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				return services.NewCheatsService(), nil
			},
		},
	}
}
