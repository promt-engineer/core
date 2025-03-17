package entities

import (
	"encoding/json"
	"time"

	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"bitbucket.org/play-workspace/base-slot-server/pkg/overlord"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type GameState struct {
	UserID         uuid.UUID `json:"user_id"`
	ExternalUserID string    `json:"external_user_id"`
	SessionToken   uuid.UUID `json:"session_token"`
	Game           string    `json:"game"`
	GameID         uuid.UUID `json:"-"`
	Integrator     string    `json:"integrator"`
	Operator       string    `json:"operator"`
	Provider       string    `json:"provider"`

	Currency           string  `json:"currency"`
	CurrencySynonym    string  `json:"currency_synonym"`
	Balance            int64   `json:"balance"`
	WagerLevels        []int64 `json:"wager_levels"`
	CurrencyMultiplier int64   `json:"currency_multiplier"`
	DefaultWager       int64   `json:"default_wager"`

	EngineInfo interface{} `json:"engine_info"`
	BootInfo   interface{} `json:"boot_info"`

	GameResults GameResults `json:"game_results"`

	IsDemo bool `json:"is_demo"`

	RTP        *int64  `json:"-"`
	Volatility *string `json:"-"`

	BuyBonus     bool `json:"buy_bonus"`
	DoubleChance bool `json:"double_chance"`
	Gamble       bool `json:"gamble"`

	AvailableRTP        []int64  `json:"available_rtp"`
	AvailableVolatility []string `json:"available_volatility"`
	OnlineVolatility    bool     `json:"online_volatility"`
	UserLocale          string   `json:"user_locale"`
	GambleDoubleUp      int64    `json:"gamble_double_up"`

	Jurisdiction string `json:"jurisdiction"`
	LobbyUrl     string `json:"lobby_url"`
	ShowCheats   bool   `json:"show_cheats"`
	LowBalance   bool   `json:"low_balance"`
	ShortLink    bool   `json:"short_link"`
	MinWager     int64  `json:"min_wager"`
}

func (gs *GameState) Compute() *GameState {
	lo.ForEach(gs.GameResults, func(item *GameResult, index int) {
		item.Compute(&gs.GambleDoubleUp)
	})

	return gs
}

func (gs *GameState) SetEngineInfo(engineInfo interface{}) *GameState {
	gs.EngineInfo = engineInfo

	return gs
}

func (gs *GameState) SetBootInfo(bootInfo interface{}) *GameState {
	gs.BootInfo = bootInfo

	return gs
}

func (gs *GameState) SetGeneratedSpin(spin engine.Spin, restoringIndexes engine.RestoringIndexes, isPFR bool, newBalance int64, roundID uuid.UUID) *HistoryRecord {
	oldBalance := newBalance - engine.TotalAward(spin) + spin.Wager()

	return gs.setGeneratedSpin(spin, restoringIndexes, isPFR, newBalance, oldBalance, roundID)
}

func (gs *GameState) SetGeneratedFreeSpin(spin engine.Spin, restoringIndexes engine.RestoringIndexes, isPFR bool, newBalance int64, roundID uuid.UUID) *HistoryRecord {
	oldBalance := newBalance
	if newBalance != 0 {
		oldBalance = newBalance - engine.TotalAward(spin)
	}

	return gs.setGeneratedSpin(spin, restoringIndexes, isPFR, newBalance, oldBalance, roundID)
}

// UpdateLastSpin is a function for gamble feature that returns history records
// with the correct start balance and end balance
func (gs *GameState) UpdateLastSpin(newSpin engine.Spin, newBalance int64) *HistoryRecord {
	if len(gs.GameResults) == 0 {
		return nil
	}

	oldRes := gs.GameResults[len(gs.GameResults)-1]
	oldBalance := newBalance - engine.TotalAwardWithGambling(newSpin) + newSpin.Wager()

	if oldRes.IsPFR {
		oldBalance -= newSpin.Wager()
	}

	oldRes.Spin = newSpin
	gs.GameResults[len(gs.GameResults)-1] = oldRes

	hr := gs.extractHistoryRecord(oldRes.Spin, oldRes.RestoringIndexes, oldRes.IsPFR, newBalance, oldBalance, oldRes.ID)

	gs.Balance = newBalance

	return hr
}

func (gs *GameState) setGeneratedSpin(spin engine.Spin, restoringIndexes engine.RestoringIndexes, isPFR bool, newBalance, oldBalance int64, roundID uuid.UUID) *HistoryRecord {
	hr := gs.extractHistoryRecord(spin, restoringIndexes, isPFR, newBalance, oldBalance, roundID)

	ngr := NewGameResult(hr.ID, spin, restoringIndexes, hr.IsPFR, gs.CurrencyMultiplier)

	if engine.GetFromContainer().HistoryHandlingType == engine.ParallelRestoring {
		gs.GameResults = append(gs.GameResults, ngr)
	} else {
		gs.GameResults = GameResults{ngr}
	}

	gs.Balance = newBalance

	return hr
}

func (gs *GameState) extractHistoryRecord(spin engine.Spin, restoringIndexes engine.RestoringIndexes, isPFR bool, newBalance, oldBalance int64, roundID uuid.UUID) *HistoryRecord {
	return &HistoryRecord{
		ID:             roundID,
		Game:           gs.Game,
		GameID:         gs.GameID,
		UserID:         gs.UserID,
		ExternalUserID: gs.ExternalUserID,
		SessionToken:   gs.SessionToken,
		Integrator:     gs.Integrator,
		Operator:       gs.Operator,
		Provider:       gs.Provider,

		Currency:     gs.Currency,
		StartBalance: oldBalance,
		EndBalance:   newBalance,
		Wager:        spin.Wager(),
		BaseAward:    spin.BaseAward(),
		BonusAward:   spin.BonusAward(),
		FinalAward:   engine.TotalAwardWithGambling(spin),

		Spin:             spin,
		RestoringIndexes: restoringIndexes,
		IsShown:          restoringIndexes.IsShown(spin),
		IsPFR:            isPFR,
		IsDemo:           gs.IsDemo,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (gs *GameState) SetRestoredSpin(hr *HistoryRecord, currencyMultiplier int64) {
	gs.GameResults = append(gs.GameResults, hr.ExtractGameResult(currencyMultiplier))

	originalWager := hr.Spin.OriginalWager()
	if lo.Contains(gs.WagerLevels, originalWager) {
		gs.DefaultWager = originalWager
	}
}

func (gs *GameState) ToWagerState() *WagerGameState {
	gs.Compute()

	return &WagerGameState{
		UserID:       gs.UserID,
		SessionToken: gs.SessionToken,
		Game:         gs.Game,
		Integrator:   gs.Integrator,

		Operator: gs.Operator,

		Currency: gs.Currency,
		Balance:  gs.Balance,

		GameResults: gs.GameResults,
	}
}

type WagerGameState struct {
	UserID       uuid.UUID `json:"user_id"`
	SessionToken uuid.UUID `json:"session_token"`
	Game         string    `json:"game"`
	Integrator   string    `json:"integrator"`
	Operator     string    `json:"operator"`

	Currency string `json:"currency"`
	Balance  int64  `json:"balance"`

	GameResults GameResults `json:"game_results"`
}

type GameResults []*GameResult

func (gr GameResults) Last() (*GameResult, bool) {
	if len(gr) == 0 {
		return nil, false
	}

	return gr[len(gr)-1], true
}

func (gr GameResults) MarshalJSON() ([]byte, error) {
	if len(gr) == 1 && engine.GetFromContainer().HistoryHandlingType == engine.SequentialRestoring {
		return json.Marshal(gr[0])
	}

	return json.Marshal([]*GameResult(gr))
}

func (gr *GameResults) Wipe() {
	*gr = []*GameResult{}
}

func GameStateFromLordState(state *overlord.InitUserStateOut) (*GameState, error) {
	cfg := engine.GetFromContainer()

	userID, err := uuid.Parse(state.UserId)
	if err != nil {
		return nil, err
	}

	sessionToken, err := uuid.Parse(state.SessionToken)
	if err != nil {
		return nil, err
	}

	gameID, err := uuid.Parse(state.GameId)
	if err != nil {
		return nil, err
	}

	defaultWager := state.DefaultWager
	wagerLevels := filterWagers(state, cfg)
	if len(wagerLevels) > 0 && !lo.Contains(wagerLevels, state.DefaultWager) {
		defaultWager = wagerLevels[0]
	}

	return &GameState{
		UserID:         userID,
		ExternalUserID: state.ExternalUserId,
		SessionToken:   sessionToken,
		Game:           state.Game,
		GameID:         gameID,
		Integrator:     state.Integrator,
		Operator:       state.Operator,
		Provider:       state.Provider,

		Currency:           state.Currency,
		CurrencySynonym:    state.CurrencySynonym,
		Balance:            state.Balance,
		WagerLevels:        wagerLevels,
		DefaultWager:       defaultWager,
		MinWager:           state.MinWager,
		CurrencyMultiplier: state.CurrencyMultiplier,

		IsDemo: state.IsDemo,

		RTP:        state.Rtp,
		Volatility: state.Volatility,

		Gamble:              state.Gamble,
		BuyBonus:            state.BuyBonus,
		DoubleChance:        state.DoubleChance,
		AvailableRTP:        state.AvailableRtp,
		AvailableVolatility: state.AvailableVolatility,
		OnlineVolatility:    state.OnlineVolatility,
		UserLocale:          state.UserLocale,
		GambleDoubleUp:      state.GambleDoubleUp,

		Jurisdiction: state.Jurisdiction,
		LobbyUrl:     state.LobbyUrl,
		ShowCheats:   state.ShowCheats,
		LowBalance:   state.LowBalance,
		ShortLink:    state.ShortLink,
	}, nil
}

func filterWagers(state *overlord.InitUserStateOut, cfg *engine.Bootstrap) []int64 {
	return lo.Filter(state.WagerLevels, func(item int64, index int) bool {
		ok := true

		if cfg.AnteBetMultiplier > 0 {
			ok = item*cfg.AnteBetMultiplier/100%10 == 0
		}

		if ok && cfg.GameMaxWager > 0 {
			ok = item <= cfg.GameMaxWager*state.CurrencyMultiplier
		}

		return ok
	})
}
