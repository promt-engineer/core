package entities

import (
	"time"

	"bitbucket.org/play-workspace/base-slot-server/pkg/overlord"
)

type FreeSpin struct {
	ID         string    `json:"id"`
	Currency   string    `json:"currency"`
	ExpireDate time.Time `json:"expire_date"`
	Value      int       `json:"value"`
	Game       string    `json:"game"`
	SpinCount  int       `json:"spin_count"`
}

func FreeSpinsFromLord(bets []*overlord.FreeBet) []*FreeSpin {
	spins := make([]*FreeSpin, 0, len(bets))

	for _, bet := range bets {
		spins = append(spins, FreeSpinFromLord(bet))
	}

	return spins
}

func FreeSpinFromLord(bet *overlord.FreeBet) *FreeSpin {
	return &FreeSpin{
		ID:         bet.Id,
		Currency:   bet.Currency,
		ExpireDate: time.UnixMilli(bet.ExpireDate),
		Value:      int(bet.Value),
		Game:       bet.Game,
		SpinCount:  int(bet.SpinCount),
	}
}
