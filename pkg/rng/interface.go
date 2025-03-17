package rng

import "time"

type Client interface {
	Rand(max uint64) (rand uint64, err error)
	RandSlice(maxSlice []uint64) (rand []uint64, err error)
	RandFloat() (float64, error)
	RandFloatSlice(count int) ([]float64, error)
}

type Config struct {
	Host              string
	Port              string
	IsSecure          bool
	MaxProcessingTime time.Duration
}
