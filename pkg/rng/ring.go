package rng

import (
	"golang.org/x/exp/constraints"
	"sync"
)

// we need general lib

type Ring[T constraints.Float | constraints.Integer] struct {
	len         int
	ring        []T
	read, write int
	mu          sync.Mutex
}

func NewRing[T constraints.Float | constraints.Integer](len int) *Ring[T] {
	if len <= 0 {
		panic("ring can not have negative size")
	}

	return &Ring[T]{len: len, ring: make([]T, len)}
}

func (r *Ring[T]) Read() (T, bool) {
	r.mu.Lock()
	defer r.mu.Unlock() // can not do it earlier

	if r.read == r.write {
		return 0, false
	}

	res := r.ring[r.read]
	r.read++
	r.read %= r.len

	return res, true
}

func (r *Ring[T]) Write(toWrite T) bool {
	r.mu.Lock()
	defer r.mu.Unlock() // can not do it earlier

	if (r.write+1)%r.len == r.read {
		return false
	}

	r.ring[r.write] = toWrite
	r.write++
	r.write %= r.len

	return true
}
