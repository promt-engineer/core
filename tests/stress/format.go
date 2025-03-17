package stress

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/facade"
	transportHTTP "bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/http"
	"bitbucket.org/play-workspace/base-slot-server/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
)

const (
	StatePath     = "core/state"
	WagerPath     = "core/wager"
	RestoringPath = "core/spin_indexes/update"

	DefaultWager      = 1000
	DefaultCurrency   = "usd"
	DefaultIntegrator = "mock"
)

type InitParams struct {
	Currency   string    `json:"currency"`
	Game       string    `json:"game"`
	Integrator string    `json:"integrator"`
	UserID     uuid.UUID `json:"user_id"`
}

type State struct {
	UserID       string `json:"user_id"`
	SessionToken string `json:"session_token"`
	Game         string `json:"game"`
	Integrator   string `json:"integrator"`

	Currency string `json:"currency"`
	Balance  int64  `json:"balance"`

	GameResult struct {
		ID               uuid.UUID   `json:"id" mapstructure:"id"`
		Spin             interface{} `json:"spin" mapstructure:"spin"`
		RestoringIndexes interface{} `json:"restoring_indexes" mapstructure:"restoring_indexes"`

		IsPFR bool `json:"is_pfr" mapstructure:"is_pfr"`
		// computed
		CanGamble bool `json:"can_gamble" mapstructure:"can_gamble"`
	} `json:"game_results"`
}

func ParseStateResponse(resp *http.Response) (*State, error) {
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return nil, fmt.Errorf("bad status code: %d, body: %s", resp.StatusCode, string(bytes))
	}

	trResp := &transportHTTP.Response{}
	if err := json.Unmarshal(bytes, &trResp); err != nil {
		return nil, err
	}

	res, err := utils.ReMarshal[State](trResp.Data)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func GenerateInitRequest(urlBase url.URL, game string) (*http.Request, error) {
	link, err := url.JoinPath(urlBase.String(), StatePath)
	if err != nil {
		return nil, err
	}

	body := facade.InitStateRequest{
		Integrator: DefaultIntegrator,
		Game:       game,

		OverlordParams: InitParams{
			Currency:   DefaultCurrency,
			Game:       game,
			Integrator: DefaultIntegrator,
			UserID:     uuid.New(),
		},
	}

	return Request(link, body)
}

func WagerRequest(urlBase url.URL, sessionToken string, wager int64, engineParams interface{}) (*http.Request, error) {
	link, err := url.JoinPath(urlBase.String(), WagerPath)
	if err != nil {
		return nil, err
	}

	body := facade.WagerRequest{
		SessionToken: sessionToken,
		Wager:        wager,
		EngineParams: engineParams,
	}

	return Request(link, body)
}

func UpdateRestoringRequest(urlBase url.URL, sessionToken string, engineParams interface{}) (*http.Request, error) {
	link, err := url.JoinPath(urlBase.String(), RestoringPath)
	if err != nil {
		return nil, err
	}

	body := facade.UpdateSpinIndexesRequest{
		SessionToken:     sessionToken,
		RestoringIndexes: engineParams,
	}

	return Request(link, body)
}

func Request(link string, body interface{}) (*http.Request, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, link, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://google.com")

	return req, nil
}
