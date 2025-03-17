package roulette

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"encoding/json"
)

type RestoringIndexes struct {
	IsShownVal bool `json:"is_shown"`
}

func (r *RestoringIndexes) IsShown(spin engine.Spin) bool {
	return r.IsShownVal
}

func (r *RestoringIndexes) Update(payload interface{}) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, r)
}
