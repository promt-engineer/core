package services

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/entities"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/errs"
	"bitbucket.org/play-workspace/base-slot-server/pkg/overlord"
	"context"
	"go.uber.org/zap"
)

type FreeSpinService struct {
	overlord overlord.Client
}

func NewFreeSpinService(overlord overlord.Client) *FreeSpinService {
	return &FreeSpinService{overlord: overlord}
}

func (srv *FreeSpinService) GetFreeSpins(ctx context.Context, sessionToken string) ([]*entities.FreeSpin, error) {
	overlordSpins, err := srv.overlord.GetAvailableFreeSpins(ctx, sessionToken)
	if err != nil {
		zap.S().Info("overlord error: ", err)

		return nil, errs.TranslateOverlordErr(err)
	}

	return entities.FreeSpinsFromLord(overlordSpins.FreeBets), nil
}

func (srv *FreeSpinService) CancelFreeSpins(ctx context.Context, sessionToken string) error {
	err := srv.overlord.CancelAvailableFreeSpins(ctx, sessionToken)
	if err != nil {
		zap.S().Info("overlord error: ", err)

		return errs.TranslateOverlordErr(err)
	}

	return nil
}

func (srv *FreeSpinService) GetFreeSpinsWithIntegratorBet(ctx context.Context, sessionToken string) (map[string][]*entities.FreeSpin, error) {
	overlordSpins, err := srv.overlord.GetAvailableFreeBetsWithIntegratorBet(ctx, sessionToken)
	if err != nil {
		zap.S().Info("overlord error: ", err)

		return nil, errs.TranslateOverlordErr(err)
	}
	mp := make(map[string][]*entities.FreeSpin)

	for key, betList := range overlordSpins.FreeBets {
		mp[key] = entities.FreeSpinsFromLord(betList.Bets)
	}

	return mp, nil
}

func (srv *FreeSpinService) CancelFreeSpinsWithIntegratorBet(
	ctx context.Context,
	sessionToken string,
	integratorBetId string,
) error {
	err := srv.overlord.CancelAvailableFreeBetsByIntegratorBet(ctx, sessionToken, integratorBetId)
	if err != nil {
		zap.S().Info("overlord error: ", err)

		return errs.TranslateOverlordErr(err)
	}

	return nil
}
