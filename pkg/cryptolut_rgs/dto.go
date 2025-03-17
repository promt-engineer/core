package cryptolut_rgs

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/overlord"
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"
)

type InitReq struct {
	Game       string      `json:"game"`
	Integrator string      `json:"integrator"`
	Params     interface{} `json:"params"`
}

type SessionReq struct {
	SessionToken string `json:"session_token"`
}

type OpenBetReq struct {
	RoundID      string `json:"round_id"`
	SessionToken string `json:"session_token"`
	Wager        string `json:"wager"`
}

type CloseBetReq struct {
	Award         string `json:"award"`
	SessionToken  string `json:"session_token"`
	TransactionID string `json:"transaction_id"`
}

type RollbackReq struct {
	SessionToken  string `json:"session_token"`
	TransactionID string `json:"transaction_id"`
}

type OpenBetResp struct {
	Balance       string `json:"balance"`
	TransactionID string `json:"transaction_id"`
}

type CloseBetResp struct {
	Balance string `json:"balance"`
}

type RollbackBetResp struct {
	Balance string `json:"balance"`
}

type StateDTOResp struct {
	Integrator         string   `json:"integrator"`
	Game               string   `json:"game"`
	Username           string   `json:"username"`
	Balance            string   `json:"balance"`
	Currency           string   `json:"currency"`
	UserID             string   `json:"user_id"`
	SessionToken       string   `json:"session_token"`
	DefaultWager       string   `json:"default_wager"`
	CurrencyMultiplier string   `json:"currency_multiplier"`
	WagerLevels        []string `json:"wager_levels"`
	CurrencySynonym    string   `json:"currency_synonym"`
	MinWager           string   `json:"min_wager"`
}

func (r *StateDTOResp) ToOverlord(gameName string) (*overlord.InitUserStateOut, error) {
	balance, err := cryptolutToEjawBalance(r.Balance)
	if err != nil {
		return nil, err
	}

	defaultWager, err := cryptolutToEjawBalance(r.DefaultWager)
	if err != nil {
		return nil, err
	}

	wagerLevels := []int64{}

	for _, wl := range r.WagerLevels {
		wlInt, err := cryptolutToEjawBalance(wl)
		if err != nil {
			return nil, err
		}

		wagerLevels = append(wagerLevels, wlInt)
	}

	currencyMultiplier, err := strconv.Atoi(r.CurrencyMultiplier)
	if err != nil {
		return nil, err
	}

	minWager, err := cryptolutToEjawBalance(r.MinWager)
	if err != nil {
		return nil, err
	}

	return &overlord.InitUserStateOut{
		UserId:             uuid.NewHash(sha256.New(), uuid.Nil, []byte(r.UserID), 4).String(),
		ExternalUserId:     r.UserID,
		Integrator:         IntegratorName,
		Operator:           IntegratorName,
		Provider:           IntegratorName,
		Game:               gameName,
		GameId:             uuid.NewHash(sha256.New(), uuid.Nil, []byte(gameName), 4).String(),
		Username:           r.Username,
		SessionToken:       r.SessionToken,
		Balance:            balance,
		Currency:           strings.ToLower(r.Currency),
		FreeBets:           []string{},
		DefaultWager:       defaultWager,
		CurrencyMultiplier: int64(currencyMultiplier),
		WagerLevels:        wagerLevels,
		IsDemo:             false,
		CurrencySynonym:    r.CurrencySynonym,
		MinWager:           minWager,
	}, nil
}

func cryptolutToEjawBalance(balance string) (int64, error) {
	f, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		return 0, err
	}

	return int64(f * FinancialDivider), nil
}

func ejawToCryptolutBalance(balance int64) string {
	return fmt.Sprintf("%.2f", float64(balance)/FinancialDivider)
}
