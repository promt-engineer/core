package entities

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

type GameResult struct {
	ID               uuid.UUID               `json:"id" mapstructure:"id"`
	Spin             engine.Spin             `json:"spin" mapstructure:"spin"`
	RestoringIndexes engine.RestoringIndexes `json:"restoring_indexes" mapstructure:"restoring_indexes"`

	IsPFR bool `json:"is_pfr" mapstructure:"is_pfr"`
	// computed
	CanGamble bool `json:"can_gamble" mapstructure:"can_gamble"`
	computed  bool

	currencyMultiplier int64
}

func NewGameResult(id uuid.UUID, spin engine.Spin, restoringIndexes engine.RestoringIndexes, isPFR bool, currencyMultiplier int64) *GameResult {
	return &GameResult{ID: id, Spin: spin, RestoringIndexes: restoringIndexes, IsPFR: isPFR, currencyMultiplier: currencyMultiplier}
}

func (gr *GameResult) GetCanGable(gambleDoubleUp int64) bool {
	gr.Compute(&gambleDoubleUp)

	return gr.CanGamble
}

func (gr *GameResult) Compute(gambleDoubleUp *int64) {
	gambles := gr.Spin.GetGamble()

	// cannot gamble if:
	// 1. the base award is 0
	baseAwardZero := gr.Spin.BaseAward() == 0

	// 2. gamble limit reached
	exceededLimit := false
	if gambleDoubleUp != nil {
		exceededLimit = gambles.Len() >= int(*gambleDoubleUp)
	}

	// 3. bonus triggered
	bonusTriggered := gr.Spin.BonusTriggered()

	// 4. gamble is collected (check restoring indexes)
	gambleCollected := gr.Spin.CanGamble(gr.RestoringIndexes)

	// 5. last gamble was lost
	var lastGambleWasLost bool
	if gambles.Len() > 0 {
		lastGambleWasLost = gambles.Last().Award == 0
	}

	gr.computed = true
	gr.CanGamble = !(baseAwardZero || exceededLimit || bonusTriggered || lastGambleWasLost) && gambleCollected
}

func (gr *GameResult) MarshalJSON() ([]byte, error) {
	if !gr.computed {
		gr.Compute(nil)
	}

	return json.Marshal(gr.View())
}

func (gr *GameResult) View() map[string]interface{} {
	view := map[string]interface{}{}

	// bad decoding rebuild on simple map
	_ = mapstructure.Decode(gr, &view)

	if engine.GetFromContainer().HistoryHandlingType == engine.SequentialRestoring {
		delete(view, "id")
	}

	if !engine.GetFromContainer().GambleAnyWinFeature {
		delete(view, "can_gamble")
	}

	return view
}
