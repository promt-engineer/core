package services

import (
	"context"
	"fmt"

	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/entities"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/errs"
	"bitbucket.org/play-workspace/base-slot-server/pkg/overlord"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"golang.org/x/exp/slog"
)

type GameFlowService struct {
	lord       overlord.Client
	boot       *engine.Bootstrap
	historySrv *HistoryService
	cheatsSrv  *CheatsService
}

func NewGameFlowService(lord overlord.Client, historySrv *HistoryService, cheatsSrv *CheatsService) *GameFlowService {
	return &GameFlowService{
		lord:       lord,
		boot:       engine.GetFromContainer(),
		historySrv: historySrv,
		cheatsSrv:  cheatsSrv,
	}
}

func (s *GameFlowService) InitGame(ctx context.Context,
	game, integrator string, lordParams interface{},
) (*entities.GameState, error) {
	lordState, err := s.lord.InitUserState(ctx, game, integrator, lordParams)
	if err != nil {
		return nil, errs.TranslateOverlordErr(err)
	}

	state, err := entities.GameStateFromLordState(lordState)
	if err != nil {
		return nil, err
	}

	return state.
		SetEngineInfo(s.boot.GetEngineInfo()).
		SetBootInfo(s.boot.GetBootInfo()), nil
}

func (s *GameFlowService) GameState(ctx context.Context, sessionToken string) (*entities.GameState, error) {
	lordState, err := s.lord.GetStateBySessionToken(ctx, sessionToken)
	if err != nil {
		return nil, errs.TranslateOverlordErr(err)
	}

	state, err := entities.GameStateFromLordState(lordState)
	if err != nil {
		return nil, err
	}

	return state.
		SetEngineInfo(s.boot.GetEngineInfo()).
		SetBootInfo(s.boot.GetBootInfo()), nil
}

func (s *GameFlowService) Wager(ctx context.Context,
	gameState *entities.GameState, freeSpinID string, wager int64, params interface{}, minWager int64) (
	*entities.GameState, *entities.HistoryRecord, error,
) {
	isPFR := freeSpinID != ""

	if !isPFR && !lo.Contains(gameState.WagerLevels, wager) {
		return nil, nil, errs.NewInternalValidationErrorFromString(errs.OneOfListError("wager", gameState.WagerLevels))
	}

	if isPFR {
		fs, err := s.findUserFreeSpin(ctx, gameState.SessionToken, freeSpinID)
		if err != nil {
			return nil, nil, err
		}

		wager = int64(fs.Value)

		if wager == 0 {
			wager = minWager
			err = s.saveDefaultWagerInFreeBetValue(ctx, gameState.SessionToken, freeSpinID, minWager)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	engCtx := s.getEngineContext(ctx, gameState, params)

	var (
		spin             engine.Spin
		indexes          engine.RestoringIndexes
		err              error
		award            int64
		roundID          = uuid.New()
		exceedMultiplier bool
	)

	const (
		maxMultiplier = 5000
		maxAttempts   = 100
	)

	attempts := 0
	for {
		if attempts >= maxAttempts {
			return nil, nil, fmt.Errorf("maximum number of generation attempts has been exceeded (%d)", maxAttempts)
		}

		spin, indexes, err = s.boot.SpinFactory.Generate(engCtx, wager, params)
		if err != nil {
			return nil, nil, err
		}

		award = engine.TotalAward(spin)

		if award/wager < int64(maxMultiplier) {
			break
		}

		slog.Warn("Win exceeds max multiplier",
			"attempt", attempts+1,
			"award", award,
			"wager", wager,
			"win_multiplier", float64(award)/float64(wager),
			"max multiplier", maxMultiplier,
			"round_id", roundID,
		)
		exceedMultiplier = true
		attempts++
	}
	if exceedMultiplier {
		slog.Warn("Given: ",
			"award", award,
			"wager", wager,
			"round_id", roundID,
		)
	}

	// engine can modify wager (for example, ante bet game or buy bonus features)
	if !isPFR && gameState.Balance < spin.Wager() {
		return nil, nil, errs.ErrNotEnoughMoney
	}

	var record *entities.HistoryRecord

	bet, err := s.lord.AtomicBet(ctx, gameState.SessionToken.String(), freeSpinID, roundID.String(), spin.Wager(), award, false)
	if err != nil {
		return nil, nil, errs.TranslateOverlordErr(err)
	}

	if isPFR {
		record = gameState.SetGeneratedFreeSpin(spin, indexes, isPFR, bet.Balance, roundID)
	} else {
		record = gameState.SetGeneratedSpin(spin, indexes, isPFR, bet.Balance, roundID)
	}

	transactionID, err := uuid.Parse(bet.TransactionId)
	if err != nil {
		// rollback
		return nil, nil, errs.ErrInternalBadData
	}

	record.SetTransactionID(transactionID)

	return gameState, record, nil
}

func (s *GameFlowService) GambleAnyWin(ctx context.Context, gameState *entities.GameState, params interface{}) (*entities.GameState, *entities.HistoryRecord, error) {
	if gameState.GambleDoubleUp == 0 {
		return nil, nil, errs.ErrLimitForGambleSetToZero
	}

	engCtx := s.getEngineContext(ctx, gameState, params)

	lgr, ok := gameState.GameResults.Last()
	if !ok {
		return nil, nil, errs.ErrHistoryRecordNotFound
	}

	if !lgr.GetCanGable(gameState.GambleDoubleUp) {
		return nil, nil, errs.ErrCanNotGamble
	}

	gamble := lgr.Spin.GetGamble()

	err := gamble.Play(s.boot.SpinFactory.GetRngClient(), lgr.Spin, params, engCtx.Cheats)
	if err != nil {
		return nil, nil, err
	}

	roundID := uuid.NewString()

	award, wager := gamble.Award(), gamble.Wager()

	bet, err := s.lord.AtomicBet(ctx, gameState.SessionToken.String(), "", roundID, wager, award, true)
	if err != nil {
		return nil, nil, errs.TranslateOverlordErr(err)
	}

	record := gameState.UpdateLastSpin(lgr.Spin, bet.Balance)

	transactionID, err := uuid.Parse(bet.TransactionId)
	if err != nil {
		// rollback
		return nil, nil, errs.ErrInternalBadData
	}

	record.SetTransactionID(transactionID)

	return gameState, record, nil
}

func (s *GameFlowService) KeepGenerating(ctx context.Context, gameState *entities.GameState, params interface{}) (*entities.GameState, *entities.HistoryRecord, error) {
	engCtx := s.getEngineContext(ctx, gameState, params)

	lgr, ok := gameState.GameResults.Last()
	if !ok {
		return nil, nil, errs.ErrHistoryRecordNotFound
	}

	oldSpin := lgr.Spin.DeepCopy()

	spin, ok, err := s.boot.SpinFactory.KeepGenerate(engCtx, params)
	if err != nil {
		return nil, nil, err
	}

	if !ok {
		return nil, nil, errs.ErrSpinGenerationCanNotBeContinued
	}

	roundID := uuid.NewString()

	award := engine.TotalAward(spin) - engine.TotalAward(oldSpin)
	if award < 0 {
		zap.S().Errorf("negative award %v", award)

		return nil, nil, errs.ErrBadDataGiven
	}

	bet, err := s.lord.AtomicBet(ctx, gameState.SessionToken.String(), "", roundID, 0, award, false)
	if err != nil {
		return nil, nil, errs.TranslateOverlordErr(err)
	}

	record := gameState.UpdateLastSpin(spin, bet.Balance)

	transactionID, err := uuid.Parse(bet.TransactionId)
	if err != nil {
		zap.S().Errorf("can not parse uuid %v", bet.TransactionId)

		return nil, nil, errs.ErrInternalBadData
	}

	record.SetTransactionID(transactionID)

	return gameState, record, nil
}

func (s *GameFlowService) findUserFreeSpin(ctx context.Context, session uuid.UUID, fsID string) (*entities.FreeSpin, error) {
	pureFS, err := s.lord.GetAvailableFreeSpins(ctx, session.String())
	if err != nil {
		return nil, errs.TranslateOverlordErr(err)
	}

	fs := entities.FreeSpinsFromLord(pureFS.FreeBets)

	item, ok := lo.Find(fs, func(item *entities.FreeSpin) bool {
		return item.ID == fsID
	})

	if !ok {
		return nil, errs.ErrWrongFreeSpinID
	}

	return item, nil
}

func (s *GameFlowService) getEngineContext(ctx context.Context, gameState *entities.GameState, params interface{}) (engCtx engine.Context) {
	engCtx = engine.Context{Context: ctx}

	return s.bound(engCtx, gameState, params)
}

func (s *GameFlowService) bound(ctx engine.Context, gameState *entities.GameState, params interface{}) engine.Context {
	lgr, ok := gameState.GameResults.Last()
	if ok {
		ctx.LastSpin = lgr.Spin
	}

	if payload, ok := s.cheatsSrv.Get(gameState.SessionToken.String()); ok {
		ctx.Cheats = payload
	}

	var userRTP *int64
	var userVolatility *string

	if gameState.RTP != nil || gameState.Volatility != nil {
		userRTP = gameState.RTP
		userVolatility = gameState.Volatility
	}

	if gameState.OnlineVolatility {
		features, err := engine.UnmarshalTo[engine.Features](params)
		if err != nil {
			zap.S().Errorf("can not parse params %v", params)
		} else if features.Volatility != "" {
			userVolatility = &features.Volatility
		}
	}

	if userRTP != nil || userVolatility != nil {
		ok := s.CheckAvailableRTPAndVolatility(userRTP, userVolatility, gameState.AvailableRTP, gameState.AvailableVolatility)

		if !ok {
			zap.S().Errorf("wrong rtp: %v or volatility: %v ", userRTP, userVolatility)
			return ctx
		}
		ctx.UserParams = &engine.UserParams{
			RTP:        userRTP,
			Volatility: userVolatility,
		}
	}

	return ctx
}

func (s *GameFlowService) CheckAvailableRTPAndVolatility(rtp *int64, volatility *string, availableRTP []int64, availableVolatility []string) bool {
	wrongRTP := true
	wrongVolatility := true

	if rtp != nil {
		for _, r := range availableRTP {
			if *rtp == r {
				wrongRTP = false
				break
			}
		}
	} else {
		wrongRTP = false
	}

	if volatility != nil {
		for _, v := range availableVolatility {
			if *volatility == v {
				wrongVolatility = false
				break
			}
		}
	} else {
		wrongVolatility = false
	}

	if wrongRTP || wrongVolatility {
		return false
	}

	return true
}

func (s *GameFlowService) saveDefaultWagerInFreeBetValue(ctx context.Context, session uuid.UUID, fsID string, value int64) error {
	err := s.lord.SaveDefaultWagerInFreeBetValue(ctx, session.String(), fsID, value)
	if err != nil {
		return errs.TranslateOverlordErr(err)
	}

	return nil
}
