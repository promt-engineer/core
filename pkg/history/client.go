package history

import (
	"context"
	"crypto/tls"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	Create(ctx context.Context, record *SpinIn) error
	Update(ctx context.Context, record *SpinIn) error

	Pagination(ctx context.Context, internalUserID uuid.UUID, game string, count int, page int) (p *GetSpinPaginationOut, err error)

	LastRecord(ctx context.Context, internalUserID uuid.UUID, game string) (*SpinOut, error)
	LastRecords(ctx context.Context, internalUserID uuid.UUID, game string) ([]*SpinOut, error)
	LastRecordByWager(ctx context.Context, internalUserID uuid.UUID, game string, wager uint64) (*SpinOut, error)
	GetByID(ctx context.Context, id uuid.UUID) (*SpinOut, error)
}

type Config struct {
	Host     string
	Port     string
	IsSecure bool
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

func newClient(host, port string, isSecure bool) (HistoryServiceClient, error) {
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

	return NewHistoryServiceClient(conn), nil
}

type client struct {
	api HistoryServiceClient
}

func (c *client) Create(ctx context.Context, record *SpinIn) error {
	_, err := c.api.CreateSpin(ctx, record)

	return err
}

func (c *client) Update(ctx context.Context, record *SpinIn) error {
	_, err := c.api.UpdateSpin(ctx, record)

	return err
}

func (c *client) Pagination(ctx context.Context, internalUserID uuid.UUID, game string, count int, page int) (p *GetSpinPaginationOut, err error) {
	return c.api.GetSpinsPagination(ctx, &GetSpinPaginationIn{
		Filter: &GetLastSpinIn{
			InternalUserId: internalUserID.String(),
			Game:           game,
		},
		Limit: uint64(count),
		Page:  uint64(page),
	})
}

func (c *client) LastRecord(ctx context.Context, internalUserID uuid.UUID, game string) (*SpinOut, error) {
	res, err := c.api.GetLastSpin(ctx, &GetLastSpinIn{
		InternalUserId: internalUserID.String(),
		Game:           game,
	})

	if err != nil {
		return nil, err
	}

	if !res.IsFound {
		return nil, ErrSpinNotFound
	}

	return res.Item, nil
}

func (c *client) LastRecords(ctx context.Context, internalUserID uuid.UUID, game string) ([]*SpinOut, error) {
	res, err := c.api.GetLastNotShownSpins(ctx, &GetLastSpinIn{
		InternalUserId: internalUserID.String(),
		Game:           game,
	})

	if err != nil {
		return nil, err
	}

	return res.Items, nil
}

func (c *client) GetByID(ctx context.Context, id uuid.UUID) (*SpinOut, error) {
	res, err := c.api.GetSpin(ctx, &GetSpinIn{
		RoundId: id.String(),
	})

	if err != nil {
		return nil, err
	}

	if !res.IsFound {
		return nil, ErrSpinNotFound
	}

	return res.Item, nil
}

func (c *client) LastRecordByWager(ctx context.Context, internalUserID uuid.UUID, game string, wager uint64) (*SpinOut, error) {
	res, err := c.api.GetLastSpinByWager(ctx, &GetLastSpinByWagerIn{
		InternalUserId: internalUserID.String(),
		Game:           game,
		Wager:          wager,
	})

	if err != nil {
		return nil, err
	}

	if !res.IsFound {
		return nil, ErrSpinNotFound
	}

	return res.Item, nil
}
