package container

import (
	"context"
	"sync"

	"github.com/sarulabs/di"
)

var (
	container di.Container
	once      sync.Once
)

func Build(ctx context.Context, wg *sync.WaitGroup, configPath, buildTag string, isDebug, isAvailableCheats bool) di.Container {
	once.Do(func() {
		builder, _ := di.NewBuilder()
		defs := []di.Def{}

		defs = append(defs, BuildServices()...)
		defs = append(defs, BuildRepositories()...)
		defs = append(defs, BuildHandlers(buildTag, isAvailableCheats)...)
		defs = append(defs, BuildGenerals(configPath, isDebug)...)
		defs = append(defs, BuildPkg()...)
		defs = append(defs, BuildTransport(ctx, wg)...)
		defs = append(defs, BuildMiddlewares()...)

		if err := builder.Add(defs...); err != nil {
			panic(err)
		}

		container = builder.Build()
	})

	return container
}

func NewBuild(ctx context.Context, wg *sync.WaitGroup, configPath string) di.Container {
	once.Do(func() {
		builder, _ := di.NewBuilder()
		defs := []di.Def{}

		defs = append(defs, BuildServices()...)
		defs = append(defs, BuildRepositories()...)
		defs = append(defs, NewBuildHandlers()...)
		defs = append(defs, NewBuildGenerals(configPath)...)
		defs = append(defs, BuildPkg()...)
		defs = append(defs, BuildTransport(ctx, wg)...)
		defs = append(defs, BuildMiddlewares()...)

		if err := builder.Add(defs...); err != nil {
			panic(err)
		}

		container = builder.Build()
	})

	return container
}
