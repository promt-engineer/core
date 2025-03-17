package cryptolut_rgs

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/overlord"
	"context"
)

type Client struct {
	cfg *Config
}

func NewClient(cfg *Config) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) InitUserState(ctx context.Context, game, integrator string, lordParams interface{}) (state *overlord.InitUserStateOut, err error) {
	req := &InitReq{Game: game, Integrator: IntegratorName, Params: lordParams}
	resp := &StateDTOResp{}

	if err := askServer(ctx, buildURL(c.cfg.BaseURL, InitPath), req, &resp)(c.cfg.WithLogs); err != nil {
		return nil, err
	}

	return resp.ToOverlord(game)
}

func (c *Client) GetStateBySessionToken(ctx context.Context, token string) (*overlord.InitUserStateOut, error) {
	req := &SessionReq{SessionToken: token}
	resp := &StateDTOResp{}

	if err := askServer(ctx, buildURL(c.cfg.BaseURL, GetStatePath), req, &resp)(c.cfg.WithLogs); err != nil {
		return nil, err
	}

	return resp.ToOverlord(TMPGameName)
}

func (c *Client) AtomicBet(ctx context.Context, sessionToken, freeBetID, roundID string, wager, award int64) (*overlord.AtomicBetOut, error) {
	openReq := &OpenBetReq{
		SessionToken: sessionToken,
		RoundID:      roundID,
		Wager:        ejawToCryptolutBalance(wager),
	}

	openResp := &OpenBetResp{}

	if err := askServer(ctx, buildURL(c.cfg.BaseURL, OpenBetPath), openReq, &openResp)(c.cfg.WithLogs); err != nil {
		return nil, err
	}

	closeReq := &CloseBetReq{
		Award:         ejawToCryptolutBalance(award),
		SessionToken:  sessionToken,
		TransactionID: openResp.TransactionID,
	}

	closeResp := &CloseBetResp{}

	if err := askServer(ctx, buildURL(c.cfg.BaseURL, CloseBetPath), closeReq, &closeResp)(c.cfg.WithLogs); err != nil {
		rollbackReq := &RollbackReq{
			SessionToken:  sessionToken,
			TransactionID: openResp.TransactionID,
		}

		rollbackResp := &RollbackBetResp{}

		if err := askServer(ctx, buildURL(c.cfg.BaseURL, RollbackBetPath), rollbackReq, rollbackResp)(c.cfg.WithLogs); err != nil {
			return nil, err
		}

		return nil, err
	}

	balance, err := cryptolutToEjawBalance(closeResp.Balance)
	if err != nil {
		return nil, err
	}

	return &overlord.AtomicBetOut{
		Balance:       balance,
		TransactionId: openResp.TransactionID,
	}, nil

}

func (c *Client) GetAvailableFreeSpins(ctx context.Context, sessionToken string) (*overlord.GetAvailableFreeBetsOut, error) {
	return &overlord.GetAvailableFreeBetsOut{}, nil
}

func (c *Client) CancelAvailableFreeSpins(ctx context.Context, sessionToken string) error {
	return nil
}
