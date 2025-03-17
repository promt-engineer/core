package facade

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/entities"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/errs"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/services"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/validator"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Facade struct {
	validationEngine *validator.Validator
	boot             *engine.Bootstrap
	gameFlowSrv      *services.GameFlowService
	freeSpinSrv      *services.FreeSpinService
	historySrv       *services.HistoryService
	cheatsSrv        *services.CheatsService
}

func NewFacade(validationEngine *validator.Validator,
	gameFlowSrv *services.GameFlowService, historySrv *services.HistoryService,
	freeSpinSrv *services.FreeSpinService, cheatsSrv *services.CheatsService) *Facade {
	return &Facade{
		validationEngine: validationEngine,
		boot:             engine.GetFromContainer(),
		gameFlowSrv:      gameFlowSrv,
		freeSpinSrv:      freeSpinSrv,
		historySrv:       historySrv,
		cheatsSrv:        cheatsSrv,
	}
}

func (facade *Facade) InitState(ctx context.Context, payload interface{}) (*entities.GameState, error) {
	req := InitStateRequest{}
	if err := parseRequest(payload, &req, facade.validationEngine); err != nil {
		return nil, err
	}

	gs, err := facade.gameFlowSrv.InitGame(ctx, req.Game, req.Integrator, req.OverlordParams)
	if err != nil {
		return nil, err
	}

	if err = facade.restoreGameState(ctx, gs); err != nil && !errors.Is(err, errs.ErrHistoryRecordNotFound) {
		return nil, err
	}

	return gs.Compute(), nil
}

func (facade *Facade) Wager(ctx context.Context, payload interface{}, metaData *entities.PlayerMetaData) (*entities.WagerGameState, error) {
	if err := facade.validatePlayerMetadata(metaData); err != nil {
		return nil, err
	}

	req := WagerRequest{}
	if err := parseRequest(payload, &req, facade.validationEngine); err != nil {
		return nil, err
	}

	gameState, err := facade.gameFlowSrv.GameState(ctx, req.SessionToken)
	if err != nil {
		return nil, err
	}

	if err = facade.restoreGameState(ctx, gameState); err != nil && !errors.Is(err, errs.ErrHistoryRecordNotFound) {
		return nil, err
	}

	if facade.boot.HistoryHandlingType == engine.SequentialRestoring {
		lr, ok := gameState.GameResults.Last()
		if ok {
			isShown := lr.RestoringIndexes.IsShown(lr.Spin)
			if !isShown {
				return nil, errs.ErrLastSpinWasNotShown
			}

			if facade.boot.ChainDependency && lr.Spin.Wager() != req.Wager {
				gameState.GameResults.Wipe()

				if err := facade.historySrv.RestoreLastSpinByWager(ctx, gameState, uint64(req.Wager)); err != nil {
					return nil, err
				}
			}
		}
	}

	gameState, record, err := facade.gameFlowSrv.Wager(ctx, gameState, req.FreeSpinID, req.Wager, req.EngineParams, gameState.MinWager)
	if err != nil {
		return nil, err
	}

	if facade.boot.HistoryHandlingType != engine.NoHistory {
		if err = facade.historySrv.Create(ctx, record, metaData); err != nil {
			zap.S().Error(err)
		}
	}

	return gameState.ToWagerState(), nil
}

func (facade *Facade) GambleAnyWin(ctx context.Context, payload interface{}, metaData *entities.PlayerMetaData) (*entities.WagerGameState, error) {
	if err := facade.validatePlayerMetadata(metaData); err != nil {
		return nil, err
	}

	if !facade.boot.GambleAnyWinFeature {
		return nil, errs.ErrGambleAnyWinWasDisabledOnServerLevel
	}

	req := GambleAnyWinRequest{}
	if err := parseRequest(payload, &req, facade.validationEngine); err != nil {
		return nil, err
	}

	gameState, err := facade.gameFlowSrv.GameState(ctx, req.SessionToken)
	if err != nil {
		return nil, err
	}

	if err = facade.restoreGameState(ctx, gameState); err != nil {
		return nil, err
	}

	gameState, record, err := facade.gameFlowSrv.GambleAnyWin(ctx, gameState, req.EngineParams)

	if err != nil {
		return nil, err
	}

	if err = facade.historySrv.UpdateRecord(ctx, record, metaData); err != nil {
		zap.S().Error(err)
	}

	return gameState.ToWagerState(), nil
}

func (facade *Facade) KeepGenerating(ctx context.Context, payload interface{}, metaData *entities.PlayerMetaData) (*entities.WagerGameState, error) {
	if err := facade.validatePlayerMetadata(metaData); err != nil {
		return nil, err
	}

	req := KeepGeneratingRequest{}
	if err := parseRequest(payload, &req, facade.validationEngine); err != nil {
		return nil, err
	}

	gameState, err := facade.gameFlowSrv.GameState(ctx, req.SessionToken)
	if err != nil {
		return nil, err
	}

	if err = facade.restoreGameState(ctx, gameState); err != nil {
		return nil, err
	}

	gameState, record, err := facade.gameFlowSrv.KeepGenerating(ctx, gameState, req.EngineParams)

	if err != nil {
		return nil, err
	}

	if err = facade.historySrv.UpdateRecord(ctx, record, metaData); err != nil {
		zap.S().Error(err)
	}

	return gameState.ToWagerState(), nil
}

func (facade *Facade) Paginate(ctx context.Context, payload interface{}) (*entities.HistoryPagination, error) {
	req := HistoryRequest{}
	if err := parseRequest(payload, &req, facade.validationEngine); err != nil {
		return nil, err
	}

	gameState, err := facade.gameFlowSrv.GameState(context.Background(), req.SessionToken)
	if err != nil {
		return nil, err
	}

	return facade.historySrv.Pagination(ctx, gameState.UserID, gameState.Game, *req.Count, *req.Page)
}

func (facade *Facade) UpdateSpinIndexes(ctx context.Context, payload interface{}, metaData *entities.PlayerMetaData) error {
	if err := facade.validatePlayerMetadata(metaData); err != nil {
		return err
	}

	req := UpdateSpinIndexesRequest{}
	if err := parseRequest(payload, &req, facade.validationEngine); err != nil {
		return err
	}

	if facade.boot.HistoryHandlingType == engine.ParallelRestoring {
		uuidValue, err := uuid.Parse(req.RecordID)
		if err != nil {
			return err
		}

		return facade.historySrv.UpdateSpinIndexes(ctx, uuidValue, req.RestoringIndexes, metaData)
	}

	if facade.boot.HistoryHandlingType == engine.SequentialRestoring {
		gs, err := facade.gameFlowSrv.GameState(ctx, req.SessionToken)
		if err != nil {
			return err
		}

		return facade.historySrv.UpdateLastSpinIndexes(ctx, gs.UserID, gs.Game, req.RestoringIndexes, metaData)
	}

	return errs.ErrUpdatingIsNotAllowed
}

func (facade *Facade) FreeSpins(ctx context.Context, payload interface{}) (*GetFreeSpinsResponse, error) {
	req := FreeSpinsRequest{}
	if err := parseRequest(payload, &req, facade.validationEngine); err != nil {
		return nil, err
	}

	if !facade.boot.FreeSpinsFeature {
		return nil, errs.ErrGameNotSupportsFreeSpins
	}

	freeSpins, err := facade.freeSpinSrv.GetFreeSpins(ctx, req.SessionToken)
	if err != nil {
		return nil, err
	}

	return &GetFreeSpinsResponse{FreeSpins: freeSpins}, nil
}

func (facade *Facade) CancelSpins(ctx context.Context, payload interface{}) error {
	req := FreeSpinsRequest{}
	if err := parseRequest(payload, &req, facade.validationEngine); err != nil {
		return err
	}

	if !facade.boot.FreeSpinsFeature {
		return errs.ErrGameNotSupportsFreeSpins
	}

	return facade.freeSpinSrv.CancelFreeSpins(ctx, req.SessionToken)
}

func (facade *Facade) FreeSpinsWithIntegratorBet(ctx context.Context, payload interface{}) (*GetFreeSpinsWithIntegratorBetResponse, error) {
	req := FreeSpinsRequest{}
	if err := parseRequest(payload, &req, facade.validationEngine); err != nil {
		return nil, err
	}

	if !facade.boot.FreeSpinsFeature {
		return nil, errs.ErrGameNotSupportsFreeSpins
	}

	freeSpins, err := facade.freeSpinSrv.GetFreeSpinsWithIntegratorBet(ctx, req.SessionToken)
	if err != nil {
		return nil, err
	}

	return &GetFreeSpinsWithIntegratorBetResponse{FreeSpins: freeSpins}, nil
}

func (facade *Facade) CancelSpinsWithIntegratorBet(ctx context.Context, payload interface{}) error {
	req := FreeSpinsWithIntegratorBetRequest{}
	if err := parseRequest(payload, &req, facade.validationEngine); err != nil {
		return err
	}

	if !facade.boot.FreeSpinsFeature {
		return errs.ErrGameNotSupportsFreeSpins
	}

	return facade.freeSpinSrv.CancelFreeSpinsWithIntegratorBet(ctx, req.SessionToken, req.IntegratorBetId)
}

func (facade *Facade) AddCheat(_ context.Context, payload interface{}) error {
	req := CheatRequest{}
	if err := parseRequest(payload, &req, facade.validationEngine); err != nil {
		return err
	}

	facade.cheatsSrv.Add(req.SessionToken, req.Payload)

	return nil
}

func (facade *Facade) validatePlayerMetadata(playerMetadata *entities.PlayerMetaData) error {
	if err := facade.validationEngine.ValidateStruct(playerMetadata); err != nil {
		return errs.NewInternalValidationError(err)
	}

	return nil
}

func (facade *Facade) restoreGameState(ctx context.Context, gs *entities.GameState) error {
	var err error

	switch facade.boot.HistoryHandlingType {
	case engine.SequentialRestoring:
		err = facade.historySrv.SequentialRestoreGameState(ctx, gs)
	case engine.ParallelRestoring:
		err = facade.historySrv.ParallelRestoreGameState(ctx, gs)
	case engine.JustSaveHistory, engine.NoHistory:
	} // it's ok to not find last record when it's a new user

	if err != nil {
		zap.S().Info("Restore problems:", err)
	}

	return err
}

func parseRequest[T any](payload interface{}, req *T, validationEngine *validator.Validator) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &req); err != nil {
		return err
	}

	if err = validationEngine.ValidateStruct(req); err != nil {
		return errs.NewInternalValidationError(err)
	}

	return nil
}
