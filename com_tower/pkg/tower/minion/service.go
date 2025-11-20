package minion

import (
	"context"

	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/slot"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/structure"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/tower"
)

type service struct {
	integration integration
	repository *repository
}

func newService(i integration, r *repository) service {
	return service{
		integration: i,
		repository: r,
	}
}

func (s service) ListTowers() []tower.Tower {
	return s.repository.ListTowers()
}

func (s service) ListStructures() structure.Structures {
	return s.repository.ListStructures()
}

func (s service) SyncTowers(towers tower.TowersPayload) {
	s.repository.SyncTowers(towers)
}

func (s service) SyncStructures(structures structure.Structures) {
	s.repository.SyncStructures(structures)
}

func (s service) CheckSlotAvailability(ctx context.Context, request slot.SlotRequest) (slot.SlotResponse, error) {
	result, err := s.integration.RequestSlotToStructure(ctx, request)
	if err != nil {
		return slot.SlotResponse{}, err
	}

	return result, nil
}
