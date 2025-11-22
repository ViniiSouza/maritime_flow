package minion

import (
	"context"
	"fmt"

	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/slot"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/structure"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/tower"
)

type service struct {
	integration integration
	repository  *repository
}

func newService(i integration, r *repository) service {
	return service{
		integration: i,
		repository:  r,
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

func (s service) CheckSlotAvailability(ctx context.Context, request slot.SlotRequest) (*slot.SlotResponse, error) {
	result, err := s.integration.RequestSlotToStructure(ctx, request)
	if err != nil {
		return nil, err
	}

	if result.State == slot.FreeSlotState {
		acquireRequest := slot.AcquireSlotRequest{
			VehicleUUID:          request.VehicleUUID,
			StructureUUID:        request.StructureUUID,
			StructureSlotRequest: request.StructureSlotRequest,
		}

		acquireResult, err := s.integration.AcquireSlotLockInTowerLeader(ctx, acquireRequest)
		if err != nil {
			return nil, err
		}

		if acquireResult.Result == slot.UnavailableAcquireSlotResultType {
			return &slot.SlotResponse{
				State: slot.InUseSlotState,
			}, nil
		}
	}

	return result, nil
}

func (s service) SendHealthCheck(ctx context.Context) error {
	return s.integration.SendHealthCheck(ctx)
}

func (s service) ReleaseSlot(ctx context.Context, msg []byte) error {

	if err := s.integration.ReleaseSlot(ctx); err != nil {
		return fmt.Errorf("failed to release slot in structure: %v", err)
	}

	if err := s.integration.ReleaseSlotLock(ctx); err != nil {
		return fmt.Errorf("failed to release slot lock in tower leader: %v", err)
	}
}
