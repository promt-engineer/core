package rng

import (
	"context"
	"go.uber.org/zap"
	"time"
)

type SimpleClient struct {
	api               RNGClient
	MaxProcessingTime time.Duration
}

func NewSimpleClient(cfg *Config) (Client, error) {
	var err error

	client := &SimpleClient{}
	client.api, err = newClient(cfg.Host, cfg.Port, cfg.IsSecure)

	if err != nil {
		return nil, err
	}

	client.MaxProcessingTime = cfg.MaxProcessingTime

	return client, nil
}

func (c *SimpleClient) RandFloat() (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MaxProcessingTime)
	defer cancel()

	in := &RandRequestFloat{Max: uint64(1)}
	resp, err := c.api.RandFloat(ctx, in)

	if err != nil {
		zap.S().Errorf("can not rand float: %v", err)

		return 0, err
	}

	return resp.Result[0], nil
}

func (c *SimpleClient) RandFloatSlice(count int) ([]float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MaxProcessingTime)
	defer cancel()

	in := &RandRequestFloat{Max: uint64(count)}
	resp, err := c.api.RandFloat(ctx, in)

	if err != nil {
		zap.S().Errorf("can not rand float slice: %v", err)

		return nil, err
	}

	return resp.Result, nil
}

func (c *SimpleClient) Rand(max uint64) (rand uint64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MaxProcessingTime)
	defer cancel()

	in := &RandRequest{Max: []uint64{max}}
	resp, err := c.api.Rand(ctx, in)

	if err != nil {
		zap.S().Errorf("can not rand : %v", err)

		return 0, err
	}

	return resp.Result[0], nil
}

func (c *SimpleClient) RandSlice(maxSlice []uint64) (rand []uint64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MaxProcessingTime)
	defer cancel()

	in := &RandRequest{Max: maxSlice}
	resp, err := c.api.Rand(ctx, in)

	if err != nil {
		zap.S().Errorf("can not rand slice: %v", err)

		return rand, err
	}

	return resp.Result, nil
}
