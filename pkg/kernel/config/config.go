package config

import (
	"fmt"
	"path/filepath"

	"bitbucket.org/play-workspace/base-slot-server/pkg/cryptolut_rgs"
	"bitbucket.org/play-workspace/base-slot-server/pkg/history"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/constants"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/services"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/http"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/websocket"
	"bitbucket.org/play-workspace/base-slot-server/pkg/overlord"
	"bitbucket.org/play-workspace/base-slot-server/pkg/rng"
	"bitbucket.org/play-workspace/gocommon/tracer"
	"github.com/spf13/viper"
)

var (
	config *Config
)

type Config struct {
	ServerConfig         *http.Config
	WebsocketConfig      *websocket.Config
	OverlordConfig       *overlord.Config
	CryptolutConfig      *cryptolut_rgs.Config
	HistoryConfig        *history.Config
	HistoryMongoDBConfig *history.MongoDBConfig
	RNGConfig            *rng.Config
	TracerConfig         *tracer.Config

	ConstantsConfig *constants.Config
	EngineConfig    *engine.Config
	SimulatorConfig *services.SimulatorConfig
}

func New(path string) (*Config, error) {
	viper.Reset()

	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	viper.AddConfigPath(filepath.Dir(abs))
	viper.SetConfigFile(filepath.Base(abs))

	return build()
}

func build() (*Config, error) {
	config = &Config{}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	serverConfig := viper.Sub("server")
	websocketConfig := viper.Sub("websocket")
	overlordConfig := viper.Sub("overlord")
	cryptolutConfig := viper.Sub("cryptolut")
	historyConfig := viper.Sub("history")
	historyMongoDBConfig := viper.Sub("historyMongoDB")
	constantsConfig := viper.Sub("game")
	rngConfig := viper.Sub("rng")
	engineConfig := viper.Sub("engine")
	tracerConfig := viper.Sub("tracer")
	simulatorConfig := viper.Sub("simulator")

	if err := parseSubConfig(serverConfig, &config.ServerConfig); err != nil {
		return nil, err
	}

	if err := parseSubConfig(constantsConfig, &config.ConstantsConfig); err != nil {
		return nil, err
	}

	if err := parseSubConfigIfNotNil(overlordConfig, &config.OverlordConfig); err != nil {
		return nil, err
	}

	if err := parseSubConfigIfNotNil(cryptolutConfig, &config.CryptolutConfig); err != nil {
		return nil, err
	}

	if err := parseSubConfigIfNotNil(historyConfig, &config.HistoryConfig); err != nil {
		return nil, err
	}

	if err := parseSubConfigIfNotNil(historyMongoDBConfig, &config.HistoryMongoDBConfig); err != nil {
		return nil, err
	}

	if err := parseSubConfig(rngConfig, &config.RNGConfig); err != nil {
		return nil, err
	}

	if err := parseSubConfig(websocketConfig, &config.WebsocketConfig); err != nil {
		return nil, err
	}

	if err := parseSubConfig(engineConfig, &config.EngineConfig); err != nil {
		return nil, err
	}

	if err := parseSubConfigIfNotNil(simulatorConfig, &config.SimulatorConfig); err != nil {
		return nil, err
	}

	if tracerConfig != nil {
		if err := tracerConfig.Unmarshal(&config.TracerConfig); err != nil {
			panic(err)
		}
	} else {
		config.TracerConfig = &tracer.Config{Disabled: true}
	}

	return config, nil
}

func parseSubConfig[T any](subConfig *viper.Viper, parseTo *T) error {
	if subConfig == nil {
		return fmt.Errorf("can not read %T config: subconfig is nil", parseTo)
	}

	if err := subConfig.Unmarshal(parseTo); err != nil {
		return err
	}

	return nil
}

func parseSubConfigIfNotNil[T any](subConfig *viper.Viper, parseTo *T) error {
	if subConfig == nil {
		return nil
	}
	if err := subConfig.Unmarshal(parseTo); err != nil {
		return err
	}

	return nil
}
