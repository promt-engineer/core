package container

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/constants"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/http/middlewares"
	"bitbucket.org/play-workspace/gocommon/tracer"
	"github.com/go-co-op/gocron"
	"github.com/sarulabs/di"
)

func BuildMiddlewares() []di.Def {
	return []di.Def{
		{
			Name: constants.HTTPCorsMiddlewareName,
			Build: func(ctn di.Container) (interface{}, error) {
				return middlewares.CORSMiddleware, nil
			},
		},
		{
			Name: constants.HTTPSessionMiddlewareName,
			Build: func(ctn di.Container) (interface{}, error) {
				return middlewares.SessionMiddleware(), nil
			},
		},
		{
			Name: constants.HTTPSessionMuMiddlewareName,
			Build: func(ctn di.Container) (interface{}, error) {
				scheduler := ctn.Get(constants.SchedulerName).(*gocron.Scheduler)

				return middlewares.SessionMu(scheduler), nil
			},
		},
		{
			Name: constants.HTTPTraceMiddlewareName,
			Build: func(ctn di.Container) (interface{}, error) {
				tr := ctn.Get(constants.TracerName).(*tracer.JaegerTracer)

				return tracer.Trace(tr), nil
			},
		},
	}
}
