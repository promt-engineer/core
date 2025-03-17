package main

import (
	"bitbucket.org/play-workspace/base-slot-server/internal/roulette"
	"bitbucket.org/play-workspace/base-slot-server/pkg/app"
)

func main() {
	application, err := app.New("config.yml", roulette.GameBootV2)
	if err != nil {
		panic(err)
	}

	if err := application.RunOrSimulate(nil); err != nil {
		panic(err)
	}
}
