package facade

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/entities"
)

type GetFreeSpinsResponse struct {
	FreeSpins []*entities.FreeSpin `json:"freespins"`
}

type GetFreeSpinsWithIntegratorBetResponse struct {
	FreeSpins map[string][]*entities.FreeSpin `json:"freespins"`
}
