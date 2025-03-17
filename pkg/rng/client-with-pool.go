package rng

import (
	"container/list"
	"context"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	poolSize = 16
	ringSize = 128
)

// test about memmory leaks.
type WithPoolClient struct {
	// think about blocking optimization
	mu                sync.Mutex
	uintPool          *list.List
	api               RNGClient
	MaxProcessingTime time.Duration
}

type Node struct {
	max  uint64
	ring *Ring[uint64]
}

func NewWithPoolClient(cfg *Config) (Client, error) {
	var err error

	client := &WithPoolClient{uintPool: list.New()}
	client.api, err = newClient(cfg.Host, cfg.Port, cfg.IsSecure)

	if err != nil {
		zap.S().Debug(err)

		return nil, err
	}

	client.MaxProcessingTime = cfg.MaxProcessingTime

	return client, nil
}

func (c *WithPoolClient) getNewPool(value uint64) ([]uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MaxProcessingTime)
	defer cancel()

	resp, err := c.api.Rand(ctx, &RandRequest{Max: sliceOfValues(value, ringSize)})
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}

func (c *WithPoolClient) getUint(max uint64) (rand uint64, err error) {
	var (
		isFound = false
		node    *Node
	)

	for i := c.uintPool.Front(); i != c.uintPool.Back(); i = i.Next() {
		if node = i.Value.(*Node); node.max == max {
			isFound = true

			c.uintPool.MoveToFront(i)

			break
		}
	}

	if isFound {
		rand, ok := node.ring.Read()
		if ok {
			return rand, nil
		}
	} else {
		if c.uintPool.Len() == poolSize {
			c.uintPool.Remove(c.uintPool.Back())
		}
		node = &Node{max: max, ring: NewRing[uint64](ringSize)}
		c.uintPool.PushFront(node)
	}

	resp, err := c.getNewPool(max)
	if err != nil {
		return 0, err
	}

	rand = resp[0]

	for i := 1; i < len(resp); i++ {
		node.ring.Write(resp[i])
	}

	return rand, nil
}

func (c *WithPoolClient) Rand(max uint64) (rand uint64, err error) {
	return c.getUint(max)
}

func (c *WithPoolClient) RandSlice(maxSlice []uint64) (rand []uint64, err error) {
	for _, max := range maxSlice {
		res, err := c.getUint(max)
		if err != nil {
			zap.S().Debug(err)

			return nil, err
		}

		rand = append(rand, res)
	}

	return rand, nil
}

// make pool.
func (c *WithPoolClient) RandFloat() (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MaxProcessingTime)
	defer cancel()

	in := &RandRequestFloat{Max: uint64(1)}
	resp, err := c.api.RandFloat(ctx, in)

	if err != nil {
		zap.S().Debug(err)

		return 0, err
	}

	return resp.Result[0], nil
}

func (c *WithPoolClient) RandFloatSlice(count int) ([]float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MaxProcessingTime)
	defer cancel()

	in := &RandRequestFloat{Max: uint64(count)}
	resp, err := c.api.RandFloat(ctx, in)

	if err != nil {
		zap.S().Debug(err)

		return nil, err
	}

	return resp.Result, nil
}
