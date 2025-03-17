package engine

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/play-workspace/base-slot-server/pkg/rng"
)

const (
	Two int64 = 2
)

type GambleConfig struct {
	DoubleUpLimit int `json:"double_up_limit"`
}

type GambleParams struct {
	GamblePick *uint64 `json:"gamble_pick"` // zero or one
}

func ParseGambleCheats(cheats any) (*uint64, error) {
	b, err := json.Marshal(cheats)
	if err != nil {
		return nil, err
	}

	var gc GambleParams
	if err := json.Unmarshal(b, &gc); err != nil {
		return nil, err
	}

	return gc.GamblePick, nil
}

func ParseAndValidateGambleParams(parameters interface{}) (gp GambleParams, err error) {
	b, err := json.Marshal(parameters)
	if err != nil {
		return gp, err
	}

	if err := json.Unmarshal(b, &gp); err != nil {
		return gp, err
	}

	return gp, ValidateGambleParams(gp)
}

func ValidateGambleParams(gp GambleParams) error {
	if gp.GamblePick != nil && *gp.GamblePick != 0 && *gp.GamblePick != 1 {
		return fmt.Errorf("pick must be 0 or 1")
	}

	return nil
}

type GambleItem struct {
	Wager        int64  `json:"wager"`
	Award        int64  `json:"award"`
	UserPick     uint64 `json:"user_pick"`     // zero or one
	ExpectedPick uint64 `json:"expected_pick"` // zero or one
}

func (g *GambleItem) isWin() bool {
	return g.UserPick == g.ExpectedPick
}

type Gamble []*GambleItem

func (g *Gamble) Play(rng rng.Client, spin Spin, params, cheats any) (err error) {
	if g.lose() {
		return fmt.Errorf("you losed previous gamble")
	}

	wager := spin.BaseAward()

	if wager == 0 {
		return fmt.Errorf("no win to gamble")
	}

	gambleParams, err := ParseAndValidateGambleParams(params)
	if err != nil {
		return err
	}

	var expectedPick *uint64

	if cheats != nil {
		expectedPick, err = ParseGambleCheats(cheats)
		if err != nil {
			return err
		}
	}

	if expectedPick == nil {
		ep, err := rng.Rand(2)
		if err != nil {
			return err
		}

		expectedPick = &ep
	}

	gi := &GambleItem{
		Wager:        wager,
		UserPick:     *gambleParams.GamblePick,
		ExpectedPick: *expectedPick,
	}

	g.compute(gi)

	return nil
}

func (g *Gamble) Len() int {
	if g == nil {
		return 0
	}

	return len(*g)
}

func (g *Gamble) Last() *GambleItem {
	return g.last()
}

func (g *Gamble) Award() int64 {
	return g.last().Award
}

func (g *Gamble) Wager() int64 {
	last := g.last()

	if last.Award == 0 {
		if len(*g) == 1 {
			return last.Wager
		}

		return g.beforeLast().Award
	}

	return last.Wager
}

func (g *Gamble) lose() bool {
	last := g.last()

	return last != nil && !last.isWin()
}

func (g *Gamble) compute(current *GambleItem) {
	last := g.last()

	if last != nil {
		current.Wager = last.Award
	}

	if current.isWin() {
		current.Award = current.Wager * Two
	}

	*g = append(*g, current)
}

func (g *Gamble) last() *GambleItem {
	if len(*g) == 0 {
		return nil
	}

	return (*g)[len(*g)-1]
}

func (g *Gamble) beforeLast() *GambleItem {
	if len(*g) < 2 {
		return (*g)[0]
	}

	return (*g)[len(*g)-2]
}
