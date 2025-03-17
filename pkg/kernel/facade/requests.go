package facade

type InitStateRequest struct {
	Integrator     string      `json:"integrator" form:"integrator"`
	Game           string      `json:"game" form:"game" validate:"game"`
	OverlordParams interface{} `json:"params" form:"params" validate:"required"`
}

type WagerRequest struct {
	SessionToken string      `json:"session_token" form:"session_token" validate:"required"`
	Wager        int64       `json:"wager" form:"wager"`
	FreeSpinID   string      `json:"freespin_id"`
	EngineParams interface{} `json:"engine_params"`
}

type GambleAnyWinRequest struct {
	SessionToken string      `json:"session_token" form:"session_token" validate:"required"`
	EngineParams interface{} `json:"engine_params"`
}

type KeepGeneratingRequest GambleAnyWinRequest

type HistoryRequest struct {
	SessionToken string `json:"session_token" form:"session_token" query:"session_token" validate:"required"`
	Page         *int   `json:"page" form:"page" query:"page" validate:"required,gt=0"`
	Count        *int   `json:"count" form:"count" query:"count" validate:"required,gt=0"`
}

type UpdateSpinIndexesRequest struct {
	SessionToken     string      `json:"session_token" form:"session_token" validate:"required"`
	RestoringIndexes interface{} `json:"restoring_indexes"  form:"restoring_indexes" validate:"required"`
	RecordID         string      `json:"record_id" form:"record_id"`
}

type FreeSpinsRequest struct {
	SessionToken string `json:"session_token" form:"session_token" query:"session_token" validate:"required"`
}

type FreeSpinsWithIntegratorBetRequest struct {
	SessionToken    string `json:"session_token" form:"session_token" query:"session_token" validate:"required"`
	IntegratorBetId string `json:"integrator_bet_id" form:"integrator_bet_id" query:"integrator_bet_id" validate:"required"`
}

type CheatRequest struct {
	SessionToken string      `json:"session_token" form:"session_token" query:"session_token" validate:"required"`
	Payload      interface{} `json:"payload" validate:"required"`
}
