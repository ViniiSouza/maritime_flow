package minion

import (
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
)

type repository struct {
	towers []types.Tower
	structures types.Structures
}

func newRepository() *repository {
	return &repository{
		towers: []types.Tower{},
		structures: types.Structures{},
	}
}

func (r *repository) ListTowers() []types.Tower {
	return r.towers
}

func (r *repository) ListStructures() types.Structures {
	return r.structures
}

func (r *repository) SyncTowers(towers types.TowersPayload) {
	r.towers = towers.Towers
}

func (r *repository) SyncStructures(structures types.Structures) {
	r.structures = structures
}
