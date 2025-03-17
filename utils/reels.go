package utils

import "github.com/samber/lo"

// MustSubstituteReels Func for simplifying transfer reels from excel to code.
func MustSubstituteReels[Old, New comparable](substitutionMap map[Old]New, toSubstitute []Old) []New {
	return lo.Map(toSubstitute, func(item Old, index int) New {
		res, ok := substitutionMap[item]
		if !ok {
			panic(item)
		}

		return res
	})
}
