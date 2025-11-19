package minion

import (
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/structure"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/tower"
)

type service struct {
	repository repository
}

func newService(r repository) service {
	return service{
		repository: r,
	}
}

func (s service) ListTowers() []tower.Tower {
	return s.repository.ListTowers()
}

func (s service) ListStructures() structure.Structures {
	return s.repository.ListStructures()
}
