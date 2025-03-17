package integrations

import (
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	DefaultWager        int
	DefaultCurrency     string
	DefaultJurisdiction string
	DefaultUserLocale   string
	MasterGame          string
	ReskinGame          string
	IntegratorMock      string
	Host                string
	Port                string
}

var (
	configOnce sync.Once
	config     = Config{}
	configErr  error
)

func NewConfig() error {
	configOnce.Do(func() {
		viper.AddConfigPath("./tests/")
		viper.SetConfigName("config")

		if configErr = viper.ReadInConfig(); configErr != nil {
			return
		}

		configErr = viper.Unmarshal(&config)
	})

	return configErr
}
