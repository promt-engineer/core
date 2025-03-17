package services

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/history"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/entities"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/errs"
	"context"
	"errors"
	"github.com/google/uuid"
)

type HistoryService struct {
	historyClient history.Client
	boot          *engine.Bootstrap
}

func NewHistoryService(historyClient history.Client) *HistoryService {
	return &HistoryService{historyClient: historyClient, boot: engine.GetFromContainer()}
}

func (s *HistoryService) Pagination(ctx context.Context,
	userID uuid.UUID, game string, count, page int) (*entities.HistoryPagination, error) {
	p, err := s.historyClient.Pagination(ctx, userID, game, count, page)
	if err != nil {
		return nil, errs.TranslateHistoryErr(err)
	}

	f := s.boot.SpinFactory

	pagination := &entities.HistoryPagination{
		Page:  int(p.Page),
		Count: int(p.Limit),
		Total: int(p.Total),
	}

	for _, item := range p.Items {
		en, err := entities.FromHistoryServiceItem(item, f)
		if err != nil {
			return nil, err
		}

		pagination.Records = append(pagination.Records, en)
	}

	return pagination, nil

}

func (s *HistoryService) Create(ctx context.Context, record *entities.HistoryRecord, metaData *entities.PlayerMetaData) error {
	spinIn, err := record.ToHistoryServiceIn(metaData)
	if err != nil {
		return err
	}

	return errs.TranslateHistoryErr(s.historyClient.Create(ctx, spinIn))
}

func (s *HistoryService) lastRecord(ctx context.Context, userID uuid.UUID, game string) (*entities.HistoryRecord, error) {
	spinOut, err := s.historyClient.LastRecord(ctx, userID, game)
	if err != nil {
		return nil, errs.TranslateHistoryErr(err)
	}

	return entities.FromHistoryServiceItem(spinOut, s.boot.SpinFactory)
}

func (s *HistoryService) lastNotShownRecords(ctx context.Context, userID uuid.UUID, game string) ([]*entities.HistoryRecord, error) {
	items, err := s.historyClient.LastRecords(ctx, userID, game)
	if err != nil {
		return nil, errs.TranslateHistoryErr(err)
	}

	hrs := []*entities.HistoryRecord{}

	for _, item := range items {
		spin, err := entities.FromHistoryServiceItem(item, s.boot.SpinFactory)
		if err != nil {
			return nil, err
		}

		hrs = append(hrs, spin)
	}

	return hrs, nil
}

func (s *HistoryService) SequentialRestoreGameState(ctx context.Context, gameState *entities.GameState) error {
	lr, err := s.lastRecord(ctx, gameState.UserID, gameState.Game)
	if err != nil {
		return err
	}

	gameState.SetRestoredSpin(lr, gameState.CurrencyMultiplier)

	return nil
}

func (s *HistoryService) ParallelRestoreGameState(ctx context.Context, gameState *entities.GameState) error {
	records, err := s.lastNotShownRecords(ctx, gameState.UserID, gameState.Game)
	if err != nil {
		return err
	}

	for _, record := range records {
		gameState.SetRestoredSpin(record, gameState.CurrencyMultiplier)
	}

	return nil
}

func (s *HistoryService) UpdateSpinIndexes(ctx context.Context, recordID uuid.UUID, restoreIndexes interface{}, metaData *entities.PlayerMetaData) error {
	spinOut, err := s.historyClient.GetByID(ctx, recordID)
	if err != nil {
		return errs.TranslateHistoryErr(err)
	}

	record, err := entities.FromHistoryServiceItem(spinOut, s.boot.SpinFactory)
	if err != nil {
		return err
	}

	return s.update(ctx, record, restoreIndexes, metaData)
}

func (s *HistoryService) UpdateLastSpinIndexes(ctx context.Context, userID uuid.UUID, game string, restoreIndexes interface{}, metaData *entities.PlayerMetaData) error {
	record, err := s.lastRecord(ctx, userID, game)

	if err != nil {
		return err
	}

	return s.update(ctx, record, restoreIndexes, metaData)
}

func (s *HistoryService) update(ctx context.Context, record *entities.HistoryRecord, restoreIndexes interface{}, metaData *entities.PlayerMetaData) error {
	if err := record.UpdateSpinIndexes(restoreIndexes, record.Spin); err != nil {
		return err
	}

	return s.UpdateRecord(ctx, record, metaData)
}

func (s *HistoryService) UpdateRecord(ctx context.Context, record *entities.HistoryRecord, metaData *entities.PlayerMetaData) error {
	spinIn, err := record.ToHistoryServiceIn(metaData)
	if err != nil {
		return err
	}

	return errs.TranslateHistoryErr(s.historyClient.Update(ctx, spinIn))
}

func (s *HistoryService) RestoreLastSpinByWager(ctx context.Context, gameState *entities.GameState, wager uint64) error {
	spin, err := s.historyClient.LastRecordByWager(ctx, gameState.UserID, gameState.Game, wager)
	if errors.Is(err, history.ErrSpinNotFound) {
		return nil
	}

	if err != nil {
		return err
	}

	r, err := entities.FromHistoryServiceItem(spin, s.boot.SpinFactory)
	if err != nil {
		return err
	}

	gameState.SetRestoredSpin(r, gameState.CurrencyMultiplier)

	return nil
}
