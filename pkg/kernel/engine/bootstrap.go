package engine

import (
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

type Bootstrap struct {
	SpinFactory SpinFactory `mapstructure:"-"`

	HTTPTransport      bool `mapstructure:"-"`
	WebsocketTransport bool `mapstructure:"-"`

	FreeSpinsFeature    bool `mapstructure:"free_spins_feature"`
	GambleAnyWinFeature bool `mapstructure:"gamble_any_win_feature"`

	AnteBetMultiplier int64 `mapstructure:"ante_bet_multiplier"` // if ante bet value 1.25, you should set as 125
	ChainDependency   bool  `mapstructure:"chain_dependency"`

	GameMaxWager int64 `mapstructure:"game_max_wager"`

	HistoryHandlingType HistoryType `mapstructure:"-"`

	EngineInfo interface{} `mapstructure:"-"`
}

func (b *Bootstrap) GetEngineInfo() interface{} {
	return b.EngineInfo
}

func (b *Bootstrap) GetBootInfo() interface{} {
	view := map[string]interface{}{}

	// bad decoding rebuild on simple map
	if err := mapstructure.Decode(b, &view); err != nil {
		zap.S().Warn(err)
	}

	return view
}

type Config struct {
	RTP               string
	Volatility        string
	BuildVersion      string // any build info
	Debug             bool
	IsCheatsAvailable bool
	MockRNG           bool
}

type HistoryType int

const (
	NoHistory HistoryType = iota
	JustSaveHistory
	SequentialRestoring
	ParallelRestoring
)
