package overlord

import (
	"context"
	"crypto/tls"
	"encoding/json"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	InitUserState(ctx context.Context, game, integrator string, lordParams interface{}) (state *InitUserStateOut, err error)
	GetStateBySessionToken(ctx context.Context, token string) (*InitUserStateOut, error)
	AtomicBet(ctx context.Context, sessionToken, freeBetID, roundID string, wager, award int64, isGamble bool) (*AtomicBetOut, error)

	GetAvailableFreeSpins(ctx context.Context, sessionToken string) (*GetAvailableFreeBetsOut, error)
	CancelAvailableFreeSpins(ctx context.Context, sessionToken string) error
	GetAvailableFreeBetsWithIntegratorBet(ctx context.Context, sessionToken string) (*GetAvailableFreeBetsWithIntegratorBetOut, error)
	CancelAvailableFreeBetsByIntegratorBet(ctx context.Context, sessionToken string, integratorBetId string) error
	SaveDefaultWagerInFreeBetValue(ctx context.Context, sessionToken string, freeBetID string, value int64) error
}

type client struct {
	api OverlordClient
}

type Config struct {
	Host     string
	Port     string
	IsSecure bool
}

type OpenBetResponse struct {
	TransactionID string
	Balance       int64
}

type CloseBetResponse struct {
	Balance int64
}

func newClient(host, port string, isSecure bool) (OverlordClient, error) {
	addr := host + ":" + port

	var (
		conn *grpc.ClientConn
		err  error
	)

	if isSecure {
		config := &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
		}

		conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	} else {
		conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if err != nil {
		zap.S().Errorf("can not dial %v: %v", addr, err)

		return nil, err
	}

	return NewOverlordClient(conn), nil
}

func NewClient(cfg *Config) (Client, error) {
	var err error

	service := &client{}
	service.api, err = newClient(cfg.Host, cfg.Port, cfg.IsSecure)

	if err != nil {
		return service, err
	}

	return service, nil
}

func (o *client) GetStateBySessionToken(ctx context.Context, token string) (*InitUserStateOut, error) {
	in := &GetStateBySessionTokenIn{SessionToken: token}

	state, err := o.api.GetStateBySessionToken(ctx, in)
	if err != nil {
		zap.S().Errorf("can't receive state: %s", err.Error())

		return state, mapError(err)
	}

	return state, nil
}

func (o *client) InitUserState(ctx context.Context,
	game, integrator string, lordParams interface{}) (state *InitUserStateOut, err error) {
	zap.S().Info("repo: InitUserState starting...")

	// Create API request
	initUserStateIn := &InitUserStateIn{
		Integrator: integrator,
		Game:       game,
	}

	initUserStateIn.Params, err = json.Marshal(lordParams)
	if err != nil {
		return nil, ErrMarshaling
	}
	zap.S().Info("initUserStateIn ", initUserStateIn)
	zap.S().Info("lordParams ", lordParams)

	res, err := o.api.InitUserState(ctx, initUserStateIn)
	if err != nil {
		zap.S().Errorf("init user state error: %s", err.Error())

		return nil, mapError(err)
	}

	return res, nil
}

func (o *client) OpenBet(ctx context.Context,
	sessionToken, roundID, currency string, value int64) (*OpenBetResponse, error) {
	p := &OpenBetIn{SessionToken: sessionToken, RoundId: roundID, Wager: value}

	bet, err := o.api.OpenBet(ctx, p)

	if err != nil {
		zap.S().Errorf("open bet error: %s", err.Error())

		return nil, mapError(err)
	}

	y := &OpenBetResponse{TransactionID: bet.TransactionId, Balance: bet.Balance}

	return y, nil
}

func (o *client) OpenFreeBet(ctx context.Context, sessionToken, freeBetID, roundID string) (*OpenBetResponse, error) {
	p := &OpenFreeBetIn{SessionToken: sessionToken, FreeBetId: freeBetID, RoundId: roundID}

	bet, err := o.api.OpenFreeBet(ctx, p)
	if err != nil {
		zap.S().Errorf("open free bet error: %s", err.Error())

		return nil, mapError(err)
	}

	y := &OpenBetResponse{TransactionID: bet.TransactionId, Balance: bet.Balance}

	return y, nil
}

func (o *client) CloseBet(ctx context.Context, transactionID, currency string, value int64) (*CloseBetResponse, error) {
	p := &CloseBetIn{TransactionId: transactionID, Award: value}
	bet, err := o.api.CloseBet(ctx, p)

	if err != nil {
		zap.S().Errorf("close bet error: %s", err.Error())

		return nil, mapError(err)
	}

	y := &CloseBetResponse{Balance: bet.Balance}

	return y, nil
}

func (o *client) GetCurrencies(ctx context.Context) ([]string, error) {
	res, err := o.api.GetAvailableCurrencies(ctx, &GetAvailableCurrenciesIn{})
	if err != nil {
		zap.S().Errorf("get available currencies error: %s", err.Error())

		return nil, mapError(err)
	}

	return res.Currencies, nil
}

func (o *client) GetIntegratorConfig(ctx context.Context, integrator, game string) (*GetIntegratorConfigOut, error) {
	res, err := o.api.GetIntegratorConfig(ctx, &GetIntegratorConfigIn{Integrator: integrator, Game: game})
	if err != nil {
		zap.S().Errorf("get integrator config error: %s", err.Error())

		return nil, mapError(err)
	}

	return res, nil
}

func (o *client) GetAvailableFreeSpins(ctx context.Context, sessionToken string) (*GetAvailableFreeBetsOut, error) {
	zap.S().Info("repo: GetAvailableFreeSpins starting...")

	in := &GetAvailableFreeBetsIn{SessionToken: sessionToken}

	freeBets, err := o.api.GetAvailableFreeBets(ctx, in)
	if err != nil {
		zap.S().Errorf("get available free spins error: %s", err.Error())

		return freeBets, mapError(err)
	}

	return freeBets, err
}

func (o *client) CancelAvailableFreeSpins(ctx context.Context, sessionToken string) error {
	in := &CancelAvailableFreeBetsIn{SessionToken: sessionToken}

	_, err := o.api.CancelAvailableFreeBets(ctx, in)
	if err != nil {
		zap.S().Errorf("cancel available free spins error: %s", err.Error())

		return mapError(err)
	}

	return err
}

func (o *client) GetAvailableFreeBetsWithIntegratorBet(
	ctx context.Context,
	sessionToken string,
) (*GetAvailableFreeBetsWithIntegratorBetOut, error) {
	in := &GetAvailableFreeBetsIn{SessionToken: sessionToken}
	data, err := o.api.GetAvailableFreeBetsWithIntegratorBet(ctx, in)
	if err != nil {
		zap.S().Errorf("cancel available free spins error: %s", err.Error())

		return data, mapError(err)
	}

	return data, err
}

func (o *client) CancelAvailableFreeBetsByIntegratorBet(
	ctx context.Context,
	sessionToken string,
	integratorBetId string,
) error {
	in := &CancelAvailableFreeBetsByIntegratorBetIn{SessionToken: sessionToken, IntegratorBetId: integratorBetId}
	_, err := o.api.CancelAvailableFreeBetsByIntegratorBet(ctx, in)
	if err != nil {
		zap.S().Errorf("cancel available free spins error: %s", err.Error())

		return mapError(err)
	}

	return err
}

func (o *client) AddFreeSpins(ctx context.Context, in *AddFreeBetIn) (*AddFreeBetOut, error) {
	return o.api.AddFreeBets(ctx, in)
}

func (o *client) CancelFreeSpins(ctx context.Context, in *CancelFreeBetIn) (out *CancelFreeBetOut, err error) {
	out, err = o.api.CancelFreeBets(ctx, in)
	if err != nil {
		zap.S().Errorf("open bet error: %s", err.Error())

		return nil, mapError(err)
	}

	return
}

func (o *client) AtomicBet(ctx context.Context, sessionToken, freeBetID, roundID string, wager, award int64, isGamble bool) (out *AtomicBetOut, err error) {
	req := &AtomicBetIn{
		SessionToken: sessionToken,
		FreeBetId:    freeBetID,
		RoundId:      roundID,
		Wager:        wager,
		Award:        award,
		IsGamble:     isGamble,
	}

	out, err = o.api.AtomicBet(ctx, req)
	if err != nil {
		zap.S().Errorf("atomic bet error: %s, data: %v", err.Error(), req)

		return nil, mapError(err)
	}

	return
}

func (o *client) SaveDefaultWagerInFreeBetValue(ctx context.Context, sessionToken string, freeBetID string, value int64) error {
	zap.S().Info("repo: SaveDefaultWagerInFreeBetValue starting...")

	in := &SaveDefaultWagerInFreeBetValueIn{SessionToken: sessionToken, Id: freeBetID, Value: value}

	_, err := o.api.SaveDefaultWagerInFreeBetValue(ctx, in)
	if err != nil {
		zap.S().Errorf("save default wager in free bet value: %s", err.Error())

		return mapError(err)
	}

	return err
}
