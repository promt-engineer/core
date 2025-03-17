package entities

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/history"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type HistoryRecord struct {
	CreatedAt time.Time `json:"created_at" mapstructure:"-"`
	UpdatedAt time.Time `json:"updated_at" mapstructure:"-"`

	ID             uuid.UUID `json:"id" mapstructure:"id"`
	GameID         uuid.UUID `json:"-" `
	Game           string    `json:"game" mapstructure:"game"`
	UserID         uuid.UUID `json:"user_id" mapstructure:"user_id"`
	ExternalUserID string    `json:"external_user_id" mapstructure:"external_user_id"`
	SessionToken   uuid.UUID `json:"session_token" mapstructure:"session_token"`
	TransactionID  uuid.UUID `json:"transaction_id" mapstructure:"transaction_id"`
	Integrator     string    `json:"integrator" mapstructure:"integrator"`
	Operator       string    `json:"operator" mapstructure:"operator"`
	Provider       string    `json:"provider" mapstructure:"provider"`

	Currency     string `json:"currency" mapstructure:"currency"`
	StartBalance int64  `json:"start_balance" mapstructure:"start_balance"`
	EndBalance   int64  `json:"end_balance" mapstructure:"end_balance"`
	Wager        int64  `json:"wager" mapstructure:"wager"`
	BaseAward    int64  `json:"base_award" mapstructure:"base_award"`
	BonusAward   int64  `json:"bonus_award" mapstructure:"bonus_award"`
	FinalAward   int64  `json:"final_award" mapstructure:"final_award"`

	Spin             engine.Spin             `json:"spin" gorm:"serializer:spin" mapstructure:"spin"`
	RestoringIndexes engine.RestoringIndexes `json:"restoring_indexes" gorm:"serializer:restoring" mapstructure:"restoring_indexes"`

	IsShown bool `json:"is_shown" mapstructure:"is_shown"`
	IsPFR   bool `json:"is_pfr" mapstructure:"is_pfr"`
	IsDemo  bool `json:"is_demo" mapstructure:"is_demo"`
}

func (hr *HistoryRecord) ToMap() map[string]interface{} {
	res := map[string]interface{}{}

	if err := mapstructure.Decode(hr, &res); err != nil {
		zap.S().Error(err)
	}

	res["created_at"] = hr.CreatedAt
	res["updated_at"] = hr.UpdatedAt

	return res
}

func (hr *HistoryRecord) UpdateSpinIndexes(newSpinIndexes interface{}, spin engine.Spin) error {
	err := hr.RestoringIndexes.Update(newSpinIndexes)
	if err != nil {
		return err
	}

	hr.IsShown = hr.RestoringIndexes.IsShown(spin)

	return nil
}

func (hr *HistoryRecord) ToHistoryServiceIn(metaData *PlayerMetaData) (*history.SpinIn, error) {
	restoring, err := json.Marshal(hr.RestoringIndexes)
	if err != nil {
		return nil, err
	}

	details, err := json.Marshal(hr.Spin)
	if err != nil {
		return nil, err
	}

	return &history.SpinIn{
		CreatedAt: timestamppb.New(hr.CreatedAt),
		UpdatedAt: timestamppb.New(hr.UpdatedAt),

		Host:      metaData.Host,
		ClientIp:  metaData.IP,
		UserAgent: metaData.UserAgent,
		Request:   metaData.Request,

		Id:             hr.ID.String(),
		GameId:         hr.GameID.String(),
		Game:           hr.Game,
		SessionToken:   hr.SessionToken.String(),
		TransactionId:  hr.TransactionID.String(),
		Integrator:     hr.Integrator,
		Operator:       hr.Operator,
		Provider:       hr.Provider,
		InternalUserId: hr.UserID.String(),
		ExternalUserId: hr.ExternalUserID,

		Currency:     hr.Currency,
		StartBalance: uint64(hr.StartBalance),
		EndBalance:   uint64(hr.EndBalance),
		Wager:        uint64(hr.Wager),
		BaseAward:    uint64(hr.BaseAward),
		BonusAward:   uint64(hr.BonusAward),
		FinalAward:   uint64(hr.FinalAward),

		RestoringIndexes: restoring,
		Details:          details,

		IsPfr:   hr.IsPFR,
		IsShown: hr.IsShown,
		IsDemo:  &hr.IsDemo,
	}, nil
}

func FromHistoryServiceItem(spin *history.SpinOut, factory engine.SpinFactory) (*HistoryRecord, error) {
	spinDetails, err := factory.UnmarshalJSONSpin(spin.Details)
	if err != nil {
		return nil, err
	}

	restoringIndexes, err := factory.UnmarshalJSONRestoringIndexes(spin.RestoringIndexes)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(spin.Id)
	if err != nil {
		return nil, err
	}

	gameID, err := uuid.Parse(spin.GameId)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(spin.InternalUserId)
	if err != nil {
		return nil, err
	}

	sessionToken, err := uuid.Parse(spin.SessionToken)
	if err != nil {
		return nil, err
	}

	transactionID, err := uuid.Parse(spin.TransactionId)
	if err != nil {
		return nil, err
	}

	return &HistoryRecord{
		CreatedAt: spin.CreatedAt.AsTime(),
		UpdatedAt: spin.UpdatedAt.AsTime(),

		ID:             id,
		GameID:         gameID,
		Game:           spin.Game,
		UserID:         userID,
		ExternalUserID: spin.ExternalUserId,
		SessionToken:   sessionToken,
		TransactionID:  transactionID,
		Integrator:     spin.Integrator,
		Operator:       spin.Operator,
		Provider:       spin.Provider,

		Currency:     spin.Currency,
		StartBalance: int64(spin.StartBalance),
		EndBalance:   int64(spin.EndBalance),
		Wager:        int64(spin.Wager),
		BaseAward:    int64(spin.BaseAward),
		BonusAward:   int64(spin.BonusAward),
		FinalAward:   int64(spin.FinalAward),

		Spin:             spinDetails,
		RestoringIndexes: restoringIndexes,

		IsShown: *spin.IsShown,
		IsPFR:   *spin.IsPfr,
		IsDemo:  *spin.IsDemo,
	}, nil
}

// TODO: think about transaction id in simple game and in gamble
func (hr *HistoryRecord) SetTransactionID(transactionID uuid.UUID) {
	hr.TransactionID = transactionID
}

func (hr *HistoryRecord) ExtractGameResult(currencyMultiplier int64) *GameResult {
	return NewGameResult(hr.ID, hr.Spin, hr.RestoringIndexes, hr.IsPFR, currencyMultiplier)
}

type HistoryPagination struct {
	Records []*HistoryRecord `json:"records"`
	Page    int              `json:"page"`
	Count   int              `json:"count"`
	Total   int              `json:"total"`
}
