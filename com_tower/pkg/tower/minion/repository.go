package minion

import (
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/tower"
)

type repository struct {
	towers []tower.Tower
	structures types.Structures
}

func newRepository() *repository {
	return &repository{
		towers: []tower.Tower{},
		structures: types.Structures{},
	}
}

func (r *repository) ListTowers() []tower.Tower {
	return r.towers
}

func (r *repository) ListStructures() types.Structures {
	return r.structures
}

func (r *repository) SyncTowers(towers tower.TowersPayload) {
	r.towers = towers.Towers
}

func (r *repository) SyncStructures(structures types.Structures) {
	r.structures = structures
}
