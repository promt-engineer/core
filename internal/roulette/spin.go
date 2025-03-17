package roulette

import "bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"

type Spin struct {
	AwardVal     int64 `json:"award"`
	WagerVal     int64 `json:"wager"`
	CurrentValue int   `json:"current_value"`
	MaxValue     int   `json:"max_value"`
}

func (s *Spin) DeepCopy() engine.Spin {
	return &Spin{
		AwardVal:     s.AwardVal,
		WagerVal:     s.WagerVal,
		CurrentValue: s.CurrentValue,
		MaxValue:     s.MaxValue,
	}
}

func (s *Spin) BaseAward() int64 {
	return s.AwardVal
}

func (s *Spin) BonusAward() int64 {
	return 0
}

func (s *Spin) OriginalWager() int64 {
	return s.WagerVal
}

func (s *Spin) Wager() int64 {
	return s.WagerVal
}

func (s *Spin) WagerNoGambles() int64 {
	return s.WagerVal
}

func (s *Spin) GetGamble() *engine.Gamble {
	return nil
}

func (s *Spin) CanGamble(_ engine.RestoringIndexes) bool {
	return false
}

func (s *Spin) BonusTriggered() bool {
	return false
}
