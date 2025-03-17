package history

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/validator"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

const SpinsCollectionName = "spins"

type Spin struct {
	ID        string    `bson:"id" json:"id" gorm:"primaryKey" csv:"id" xlsx:"ID"`
	CreatedAt time.Time `bson:"created_at" json:"created_at" csv:"created_at" xlsx:"Created At"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at" csv:"updated_at" xlsx:"Updated At"`

	Day       time.Time `bson:"day" json:"day" csv:"day" xlsx:"Day"`
	Country   *string   `bson:"country" json:"country" gorm:"<-:create"`
	Host      string    `bson:"host" json:"host" csv:"host" xlsx:"Host" gorm:"<-:create" validate:"required,url"`
	ClientIP  string    `bson:"client_ip" json:"client_ip" csv:"client_ip" gorm:"<-:create" validate:"required,ip_addr"`
	UserAgent string    `bson:"user_agent" json:"user_agent" csv:"user_agent" gorm:"<-:create" validate:"required"`

	Request       bson.M `bson:"request" json:"request" csv:"request" gorm:"<-:create" validate:"required"`
	GameID        string `bson:"game_id" json:"game_id" gorm:"primaryKey" csv:"game_id" xlsx:"-"`
	Game          string `bson:"game" json:"game" csv:"game" xlsx:"Game;Name" validate:"required"`
	SessionToken  string `bson:"session_token" json:"session_token" csv:"session_token" xlsx:"Session Token"`
	TransactionID string `bson:"transaction_id" json:"transaction_id" csv:"transaction_id" xlsx:"Transaction ID"`

	Integrator string `bson:"integrator" json:"integrator" csv:"integrator" xlsx:"Integrator" validate:"required"`
	Operator   string `bson:"operator" json:"operator" csv:"operator" xlsx:"Operator" validate:"required"`
	Provider   string `bson:"provider" json:"provider" csv:"provider" xlsx:"Provider" validate:"required"`

	InternalUserID string `bson:"internal_user_id" json:"internal_user_id" csv:"internal_user_id" xlsx:"Internal User ID"`
	ExternalUserID string `bson:"external_user_id" json:"external_user_id" csv:"external_user_id" xlsx:"External User ID" validate:"required"`
	Currency       string `bson:"currency" json:"currency" csv:"currency" xlsx:"Currency" validate:"required"`

	StartBalance float64 `bson:"start_balance" json:"start_balance" csv:"start_balance" xlsx:"Start Balance"`
	EndBalance   float64 `bson:"end_balance" json:"end_balance" csv:"end_balance" xlsx:"End Balance"`
	Wager        float64 `bson:"wager" json:"wager" csv:"wager" xlsx:"Wager"`
	BaseAward    float64 `bson:"base_award" json:"base_award" csv:"base_award" xlsx:"Base Award"`
	BonusAward   float64 `bson:"bonus_award" json:"bonus_award" csv:"bonus_award" xlsx:"Bonus Award"`
	FinalAward   float64 `bson:"final_award" json:"final_award" csv:"final_award" xlsx:"Final Award"`

	Details          bson.M `bson:"details" json:"details" gorm:"serializer:json" csv:"-" swaggertype:"string" xlsx:"-" validate:"required"`
	RestoringIndexes bson.M `bson:"restoring_indexes" json:"restoring_indexes" gorm:"serializer:json" csv:"-" swaggertype:"string" xlsx:"-" validate:"required"`

	IsShown bool `bson:"is_shown" json:"is_shown" csv:"-" xlsx:"isShown"`
	IsPFR   bool `bson:"is_pfr" json:"is_pfr" csv:"is_pfr" xlsx:"isPFR"`
	IsDemo  bool `bson:"is_demo" json:"is_demo" csv:"is_demo" xlsx:"isDemo"`
}

func (s *Spin) ToAPIResponse() *SpinOut {
	country := ""

	if s.Country != nil {
		country = *s.Country
	}

	spinOut := &SpinOut{
		CreatedAt: timestamppb.New(s.CreatedAt),
		UpdatedAt: timestamppb.New(s.UpdatedAt),

		Country:   country,
		Host:      s.Host,
		ClientIp:  s.ClientIP,
		UserAgent: s.UserAgent,
		Request:   nil,

		Id:             s.ID,
		GameId:         s.GameID,
		Game:           s.Game,
		SessionToken:   s.SessionToken,
		TransactionId:  s.TransactionID,
		Integrator:     s.Integrator,
		Operator:       s.Operator,
		Provider:       s.Provider,
		InternalUserId: s.InternalUserID,
		ExternalUserId: s.ExternalUserID,

		Currency:     s.Currency,
		StartBalance: uint64(s.StartBalance),
		EndBalance:   uint64(s.EndBalance),
		Wager:        uint64(s.Wager),
		BaseAward:    uint64(s.BaseAward),
		BonusAward:   uint64(s.BonusAward),
		FinalAward:   uint64(s.FinalAward),

		RestoringIndexes: nil,
		Details:          nil,
		IsPfr:            &s.IsPFR,
		IsShown:          &s.IsShown,
		IsDemo:           &s.IsDemo,
	}

	var err error

	spinOut.Request, err = json.Marshal(s.Request)
	if err != nil {
		zap.S().Error(err)
	}
	spinOut.RestoringIndexes, err = json.Marshal(s.RestoringIndexes)
	if err != nil {
		zap.S().Error(err)
	}
	spinOut.Details, err = json.Marshal(s.Details)
	if err != nil {
		zap.S().Error(err)
	}

	return spinOut
}

func (s *Spin) validateSpin(validatorEngine *validator.Validator) error {
	verr := validatorEngine.ValidateStruct(s)

	var err error

	for _, taggedError := range validator.CheckValidationErrors(verr) {
		err = errors.Join(err, taggedError.Err)
	}

	if len(s.ID) != len(uuid.Nil.String()) || s.ID == uuid.Nil.String() {
		err = errors.Join(err, fmt.Errorf("field id is required"))
	}

	if len(s.GameID) != len(uuid.Nil.String()) || s.GameID == uuid.Nil.String() {
		err = errors.Join(err, fmt.Errorf("field game_id is required"))
	}

	if len(s.SessionToken) != len(uuid.Nil.String()) || s.SessionToken == uuid.Nil.String() {
		err = errors.Join(err, fmt.Errorf("field session_token is required"))
	}

	if len(s.TransactionID) != len(uuid.Nil.String()) || s.TransactionID == uuid.Nil.String() {
		err = errors.Join(err, fmt.Errorf("field transaction_id is required"))
	}

	if len(s.InternalUserID) != len(uuid.Nil.String()) || s.InternalUserID == uuid.Nil.String() {
		err = errors.Join(err, fmt.Errorf("field internal_user_id is required"))
	}

	return err
}

func spinIn2Spin(in *SpinIn) (*Spin, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, err
	}

	gameID, err := uuid.Parse(in.GameId)
	if err != nil {
		return nil, err
	}

	sessionToken, err := uuid.Parse(in.SessionToken)
	if err != nil {
		return nil, err
	}

	internalUserID, err := uuid.Parse(in.InternalUserId)
	if err != nil {
		return nil, err
	}

	transactionID, err := uuid.Parse(in.TransactionId)
	if err != nil {
		return nil, err
	}

	if in.IsDemo == nil {
		return nil, ErrIsDemoRequiredField
	}

	spin := &Spin{
		CreatedAt: in.CreatedAt.AsTime(),
		UpdatedAt: in.UpdatedAt.AsTime(),

		Host:      in.Host,
		ClientIP:  in.ClientIp,
		UserAgent: in.UserAgent,
		Request:   bson.M{},

		ID:             id.String(),
		GameID:         gameID.String(),
		Game:           in.Game,
		SessionToken:   sessionToken.String(),
		TransactionID:  transactionID.String(),
		Integrator:     in.Integrator,
		Operator:       in.Operator,
		Provider:       in.Provider,
		InternalUserID: internalUserID.String(),
		ExternalUserID: in.ExternalUserId,
		Currency:       in.Currency,

		StartBalance: float64(in.StartBalance),
		EndBalance:   float64(in.EndBalance),
		Wager:        float64(in.Wager),
		BaseAward:    float64(in.BaseAward),
		BonusAward:   float64(in.BonusAward),
		FinalAward:   float64(in.FinalAward),

		Details:          bson.M{},
		RestoringIndexes: bson.M{},

		IsShown: in.IsShown,
		IsPFR:   in.IsPfr,
		IsDemo:  *in.IsDemo,
	}

	err = json.Unmarshal(in.Request, &spin.Request)
	if err != nil {
		zap.S().Error(err)
	}
	err = json.Unmarshal(in.RestoringIndexes, &spin.RestoringIndexes)
	if err != nil {
		zap.S().Error(err)
	}
	err = json.Unmarshal(in.Details, &spin.Details)
	if err != nil {
		zap.S().Error(err)
	}

	return spin, nil
}
