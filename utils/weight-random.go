package utils

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/rng"
	"errors"
	"golang.org/x/exp/constraints"
	"math"
	"sort"
)

// NewChooserFromMap is a helper function for creating Chooser from the map
// where the key is the value and the value is the weight.
func NewChooserFromMap[T comparable, W constraints.Integer](rand rng.Client, config map[T]W) (*Chooser[T, W], error) {
	choices := make([]Choice[T, W], 0, len(config))
	for value, weight := range config {
		choices = append(choices, NewChoice(value, weight))
	}

	chooser, err := NewChooser(rand, choices...)
	if err != nil {
		return nil, err
	}

	return chooser, nil
}

// Choice is a generic wrapper that can be used to add weights for any item.
type Choice[T any, W constraints.Integer] struct {
	Item   T
	Weight W
}

// NewChoice creates a new Choice with specified item and weight.
func NewChoice[T any, W constraints.Integer](item T, weight W) Choice[T, W] {
	return Choice[T, W]{Item: item, Weight: weight}
}

// A Chooser caches many possible Choices in a structure designed to improve
// performance on repeated calls for weighted random selection.
type Chooser[T any, W constraints.Integer] struct {
	data    []Choice[T, W]
	weights []uint64
	max     uint64
	rng     rng.Client
}

// NewChooser initializes a new Chooser for picking from the provided choices.
func NewChooser[T any, W constraints.Integer](rng rng.Client, choices ...Choice[T, W]) (*Chooser[T, W], error) {
	sort.Slice(choices, func(i, j int) bool {
		return choices[i].Weight < choices[j].Weight
	})

	weights := make([]uint64, len(choices))
	var runningTotal uint64 = 0
	for i, c := range choices {
		weight := uint64(c.Weight)
		if weight < 0 {
			continue // ignore negative weights, can never be picked
		}

		if (math.MaxUint64 - runningTotal) <= weight {
			return nil, errWeightOverflow
		}
		runningTotal += weight
		weights[i] = runningTotal
	}

	if runningTotal < 1 {
		return nil, errNoValidChoices
	}

	return &Chooser[T, W]{data: choices, weights: weights, max: runningTotal, rng: rng}, nil
}

// Possible errors returned by NewChooser, preventing the creation of a Chooser
// with unsafe runtime states.
var (
	// If the sum of provided Choice weights exceed the maximum integer value
	// for the current platform (e.g. math.MaxInt32 or math.MaxInt64), then
	// the internal running total will overflow, resulting in an imbalanced
	// distribution generating improper results.
	errWeightOverflow = errors.New("sum of Choice Weights exceeds max int")
	// If there are no Choices available to the Chooser with a weight >= 1,
	// there are no valid choices and Pick would produce a runtime panic.
	errNoValidChoices = errors.New("zero Choices with Weight >= 1")
)

func (c Chooser[T, W]) MultiPick(count int) ([]T, error) {
	values := make([]uint64, count)

	for i := 0; i < count; i++ {
		values[i] = c.max
	}

	rs, err := c.rng.RandSlice(values)
	if err != nil {
		return nil, err
	}

	result := make([]T, count)

	for i, r := range rs {
		result[i] = c.data[searchInts(c.weights, r+1)].Item
	}

	return result, nil
}

// Pick returns a single weighted random Choice.Item from the Chooser.
//
// Utilizes global rand as the source of randomness.
func (c Chooser[T, W]) Pick() (T, error) {
	r, err := c.rng.Rand(c.max)
	if err != nil {
		return Choice[T, W]{}.Item, err
	}
	i := searchInts(c.weights, r+1)
	return c.data[i].Item, nil
}

// The standard library sort.SearchInts() just wraps the generic sort.Search()
// function, which takes a function closure to determine truthfulness. However,
// since this function is utilized within a for loop, it cannot currently be
// properly inlined by the compiler, resulting in non-trivial performance
// overhead.
//
// Thus, this is essentially manually inlined version.  In our use case here, it
// results in a up to ~33% overall throughput increase for Pick().
func searchInts(a []uint64, x uint64) uint64 {
	// Possible further future optimization for searchInts via SIMD if we want
	// to write some Go assembly code: http://0x80.pl/articles/simd-search.html
	var i, j uint64 = 0, uint64(len(a))
	for i < j {
		h := (i + j) >> 1 // avoid overflow when computing h
		if a[h] < x {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
