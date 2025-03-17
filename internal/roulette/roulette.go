package roulette

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine/utils/volatility"
	"bitbucket.org/play-workspace/base-slot-server/pkg/rng"
	"encoding/json"
)

const (
	max_       = 37
	payLine    = 18
	multiplier = 2
)

func GameBootV2(rand rng.Client, vol volatility.Type, rtp float64) *engine.Bootstrap {
	return &engine.Bootstrap{
		SpinFactory: NewSpinFactory(rand, vol, rtp),

		HTTPTransport: true,

		HistoryHandlingType: engine.SequentialRestoring,
	}
}

type SpinFactory struct {
	rand rng.Client

	vol volatility.Type
	rtp float64
}

func (s *SpinFactory) UnmarshalJSONRestoringIndexes(bytes []byte) (engine.RestoringIndexes, error) {
	ri := RestoringIndexes{}

	if err := json.Unmarshal(bytes, &ri); err != nil {
		return nil, err
	}

	return &ri, nil
}

func (s *SpinFactory) UnmarshalJSONSpin(bytes []byte) (engine.Spin, error) {
	spin := Spin{}
	err := json.Unmarshal(bytes, &spin)

	return &spin, err
}

func NewSpinFactory(rand rng.Client, vol volatility.Type, rtp float64) *SpinFactory {
	return &SpinFactory{
		rand: rand,
		vol:  vol,
		rtp:  rtp,
	}
}

func (s *SpinFactory) Generate(_ engine.Context, wager int64, _ interface{}) (engine.Spin, engine.RestoringIndexes, error) {
	randValue, err := s.rand.Rand(max_)
	if err != nil {
		// TODO: translate error
		return nil, nil, err
	}

	spin := Spin{WagerVal: wager, CurrentValue: int(randValue), MaxValue: max_}
	if isWon := spin.CurrentValue < payLine; isWon {
		spin.AwardVal = wager * multiplier
	}

	return &spin, &RestoringIndexes{}, nil
}

func (s *SpinFactory) KeepGenerate(ctx engine.Context, _ interface{}) (engine.Spin, bool, error) {
	return ctx.LastSpin, false, nil
}

func (s *SpinFactory) GetRngClient() rng.Client {
	return s.rand
}
