package rng

import (
	rnd "crypto/rand"
	"math/big"
)

type MockClient struct{}

func NewMockClient(cfg *Config) (Client, error) {
	return &MockClient{}, nil
}

func (c *MockClient) Rand(max uint64) (rand uint64, err error) {
	res, err := rnd.Int(rnd.Reader, big.NewInt(int64(max)))

	if err != nil {
		return 0, err
	}

	return res.Uint64(), nil
}

func (c *MockClient) RandSlice(maxSlice []uint64) (rand []uint64, err error) {
	rand = make([]uint64, 0, len(maxSlice))

	for _, max := range maxSlice {
		res, err := rnd.Int(rnd.Reader, big.NewInt(int64(max)))

		if err != nil {
			return nil, err
		}

		rand = append(rand, res.Uint64())
	}

	return rand, nil
}

func (c *MockClient) RandFloat() (float64, error) {
	rand, err := c.Rand(1 << 53)
	if err != nil {
		return 0, err
	}

	return float64(rand) / (1 << 53), nil
}

func (c *MockClient) RandFloatSlice(count int) (rand []float64, err error) {
	for i := 0; i < count; i++ {
		res, err := c.RandFloat()

		if err != nil {
			return nil, err
		}

		rand = append(rand, res)
	}

	return rand, nil
}
