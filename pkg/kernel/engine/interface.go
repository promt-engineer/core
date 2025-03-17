package engine

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/rng"
	"context"
)

type Context struct {
	context.Context
	Cheats     interface{}
	LastSpin   Spin
	UserParams *UserParams
}

type UserParams struct {
	RTP        *int64
	Volatility *string
}

type SpinFactory interface {
	/*
	*	Generate method serves the purpose of generating a single request for a spin
	*   This method can modify the wager that passes into it
	*	For example: the game presents a possibility to buy a bonus, but the price for this bonus can be different
	* 	for each game, so we can modify the wager in this method depending on the game rules.
	 */
	Generate(ctx Context, wager int64, parameters interface{}) (Spin, RestoringIndexes, error)

	/*
	*	KeepGenerate method serves the purpose of generating multiple requests for spins
	*   For example: in the game user has won a bonus game and must decide which type of bonus he wishes to play
	* 	In DB it will be stored as the same spin overlord will receive transaction: {wager = 0; award = newAward - oldAward}
	*   Spin - new spin, you may mutate existing spin or create another one, anyway base slot will use DeepCopy() of spin.
	* 	bool - false if according to the game rules spin generation con not be continued.
	*  	error - for technical error like serialization, rng etc.
	 */
	KeepGenerate(ctx Context, parameters interface{}) (Spin, bool, error)

	UnmarshalJSONSpin(bytes []byte) (Spin, error)
	UnmarshalJSONRestoringIndexes(bytes []byte) (RestoringIndexes, error)

	GetRngClient() rng.Client
}

type Generate func(wager int64, parameters interface{}) (Spin, RestoringIndexes, error)

type Spin interface {
	BaseAward() int64
	BonusAward() int64

	// OriginalWager returns the original wager before modifying it.
	// For example, in the case of buying a bonus or playing the ante bet.
	// Used to restoring the right wager when we play the ante bet or buy bonus.
	OriginalWager() int64
	// Wager returns the wager that must take from the user balance.
	Wager() int64

	DeepCopy() Spin

	BonusTriggered() bool

	GetGamble() *Gamble

	CanGamble(restoringIndexes RestoringIndexes) bool // only logical issues, for example, no gamble in bonus game or gamble is collected
}

type RestoringIndexes interface {
	IsShown(spin Spin) bool
	Update(payload interface{}) error
}

func TotalAwardWithGambling(spin Spin) int64 {
	gambles := spin.GetGamble()

	if gambles.Len() == 0 {
		return spin.BaseAward() + spin.BonusAward()
	}

	return gambles.Award()
}

func TotalAward(spin Spin) int64 {
	return spin.BaseAward() + spin.BonusAward()
}
