package minion

import (
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/structure"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/tower"
)

type repository struct {
	towers []tower.Tower
	structures structure.Structures
}

func newRepository() repository {
	return repository{
		towers: []tower.Tower{},
		structures: structure.Structures{},
	}
}

func (r repository) ListTowers() []tower.Tower {
	return r.towers
}

func (r repository) ListStructures() structure.Structures {
	return r.structures
}
