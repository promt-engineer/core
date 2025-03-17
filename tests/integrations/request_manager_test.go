package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type Hook func(status int, content []byte)

type RequestManager struct {
	generalState GeneralState
	stateHooks   []Hook
}

type GeneralState struct {
	UserID, SessionToken string
}

type StateRequest struct {
	Game       string      `json:"game"`
	Integrator string      `json:"integrator"`
	Params     StateParams `json:"params"`
}

type StateResponse struct {
	SessionToken string `json:"session_token"`
}

type StateParams struct {
	Currency     string `json:"currency"`
	Game         string `json:"game"`
	Jurisdiction string `json:"jurisdiction"`
	Integrator   string `json:"integrator"`
	UserID       string `json:"user_id"`
	Userlocale   string `json:"userlocale"`
}

type WagerRequest struct {
	SessionToken string `json:"session_token"`
	Wager        int    `json:"wager"`
	FreeSpinID   string `json:"free_spin_id"`
}

type HistoryRequest struct {
	SessionToken string
	Page, Count  *int
}

type HistoryResponse struct {
	Total, CurrentPage, Count int
}

type SpinIndexesRequest struct {
	SessionToken   string `json:"session_token"`
	SpinIndexBase  *int   `json:"spin_index_base"`
	SpinIndexBonus *int   `json:"spin_index_bonus"`
}

func (m HistoryRequest) String() string {
	buf := bytes.Buffer{}
	buf.WriteString("?")

	if m.SessionToken != "" {
		buf.WriteString(fmt.Sprintf("session_token=%v&", m.SessionToken))
	}

	if m.Page != nil {
		buf.WriteString(fmt.Sprintf("page=%v&", *m.Page))
	}

	if m.Count != nil {
		buf.WriteString(fmt.Sprintf("count=%v&", *m.Count))
	}

	str := buf.String()

	return str[:len(str)-1]
}

func NewRequestManager() *RequestManager {
	rm := &RequestManager{}
	rm.OnStateRequest(rm.UpdateSessionIDFromState)
	rm.GenerateNewUserID()

	return rm
}

func (rm *RequestManager) GetUserID() string {
	if rm.generalState.UserID == "" {
		rm.generalState.UserID = uuid.NewString()
	}

	return rm.generalState.UserID
}

func (rm *RequestManager) GenerateNewUserID() {
	rm.generalState.UserID = uuid.NewString()
	rm.GenerateNewSessionToken()
}

func (rm *RequestManager) GetWrongSessionToken() string {
	return uuid.NewString()
}

func (rm *RequestManager) GetSessionToken() string {
	if rm.generalState.SessionToken == "" {
		rm.SendStateRequest(rm.DefaultStateRequest())
	}

	return rm.generalState.SessionToken
}

func (rm *RequestManager) GetSessionTokenPure(body interface{}) string {
	status, content := SendRequest(http.MethodPost, StatePath, body)

	if status != http.StatusOK {
		smartPanic(string(content))
	}

	return rm.ParseStateResponse(content).SessionToken
}

func (rm *RequestManager) GenerateNewSessionToken() {
	rm.SendStateRequest(rm.DefaultStateRequest())
}

func (rm *RequestManager) DefaultStateRequest() StateRequest {
	params := StateParams{
		Currency:     config.DefaultCurrency,
		Game:         config.MasterGame,
		Jurisdiction: config.DefaultJurisdiction,
		Integrator:   config.IntegratorMock,
		UserID:       rm.GetUserID(),
		Userlocale:   config.DefaultUserLocale,
	}

	return StateRequest{
		Game:       config.MasterGame,
		Integrator: config.IntegratorMock,
		Params:     params,
	}
}

func (rm *RequestManager) DefaultWagerRequest() WagerRequest {
	return WagerRequest{
		Wager:        config.DefaultWager,
		SessionToken: manager.GetSessionToken(),
		FreeSpinID:   "",
	}
}

func (rm *RequestManager) DefaultHistoryRequest() HistoryRequest {
	one := 1

	return HistoryRequest{
		SessionToken: manager.GetSessionToken(),
		Page:         &one,
		Count:        &one,
	}
}

func (rm *RequestManager) DefaultSpinIndexesRequest() SpinIndexesRequest {
	zero := 0

	return SpinIndexesRequest{
		SessionToken:   manager.GetSessionToken(),
		SpinIndexBase:  &zero,
		SpinIndexBonus: &zero,
	}
}

func (rm *RequestManager) SendStateRequest(body interface{}) (status int, content []byte) {
	status, content = SendRequest(http.MethodPost, StatePath, body)
	rm.runStateHooks(status, content)

	return status, content
}

func (rm *RequestManager) OnStateRequest(hooks ...Hook) {
	rm.stateHooks = append(rm.stateHooks, hooks...)
}

func (rm *RequestManager) runStateHooks(status int, content []byte) {
	for _, hook := range rm.stateHooks {
		hook(status, content)
	}
}

func (rm *RequestManager) UpdateSessionIDFromState(status int, content []byte) {
	if status == http.StatusOK {
		resp := rm.ParseStateResponse(content)
		rm.generalState.SessionToken = resp.SessionToken
	}
}

func (rm *RequestManager) ParseStateResponse(content []byte) StateResponse {
	stateResp := StateResponse{}
	resp := Response{Data: &stateResp}
	err := json.Unmarshal(content, &resp)

	if err != nil {
		smartPanic(err)
	}

	return stateResp
}

func (rm *RequestManager) ParseHistoryResponse(content []byte) HistoryResponse {
	historyResp := HistoryResponse{}
	resp := Response{Data: &historyResp}
	err := json.Unmarshal(content, &resp)

	if err != nil {
		smartPanic(err)
	}

	return historyResp
}
