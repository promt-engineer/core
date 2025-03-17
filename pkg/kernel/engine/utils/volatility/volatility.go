package volatility

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/rng"
	"fmt"
	"github.com/samber/lo"
)

type Type string

const (
	LowType    Type = "low"
	MediumType Type = "medium"
	HighType   Type = "high"
)

type Low struct{}

func (m Low) Name() Type {
	return LowType
}

type Medium struct{}

func (m Medium) Name() Type {
	return MediumType
}

type High struct{}

func (m High) Name() Type {
	return HighType
}

var AvailableVolTypes = []Type{LowType, MediumType, HighType}

func VolFromStr(str string) (Type, error) {
	vol := Type(str)
	switch vol {
	case LowType, MediumType, HighType:
		return vol, nil
	}

	return "", fmt.Errorf("invalid volatility: expect one of %v, got %s", AvailableVolTypes, str)
}

type ConfigMap[T any] map[Type]*T

type Volatility[T any] interface {
	Name() Type
	Config(rand rng.Client, rtp float64) *T
}

func NewVolatilityMap[T any](rand rng.Client, rtp float64, arr ...Volatility[T]) ConfigMap[T] {
	return lo.SliceToMap(arr, func(item Volatility[T]) (Type, *T) {
		return item.Name(), item.Config(rand, rtp)
	})
}
