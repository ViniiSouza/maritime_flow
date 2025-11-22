package minion

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
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

func (s service) ListTowers() []types.Tower {
	return s.repository.ListTowers()
}

func (s service) ListStructures() types.Structures {
	return s.repository.ListStructures()
}

func (s service) SyncTowers(towers types.TowersPayload) {
	s.repository.SyncTowers(towers)
}

func (s service) SyncStructures(structures types.Structures) {
	s.repository.SyncStructures(structures)
}

func (s service) CheckSlotAvailability(ctx context.Context, request types.SlotRequest) (*types.SlotResponse, error) {
	result, err := s.integration.RequestSlotToStructure(ctx, request)
	if err != nil {
		return nil, err
	}

	if result.State == types.FreeSlotState {
		acquireRequest := types.AcquireSlotRequest{
			VehicleUUID:          request.VehicleUUID,
			StructureUUID:        request.StructureUUID,
			StructureSlotRequest: request.StructureSlotRequest,
		}

		acquireResult, err := s.integration.AcquireSlotLockInTowerLeader(ctx, acquireRequest)
		if err != nil {
			return nil, err
		}

		if acquireResult.Result == types.UnavailableAcquireSlotResultType {
			return &types.SlotResponse{
				State: types.InUseSlotState,
			}, nil
		}
	}

	return result, nil
}

func (s service) SendHealthCheck(ctx context.Context) error {
	return s.integration.SendHealthCheck(ctx)
}

func (s service) ReleaseSlot(ctx context.Context, data []byte) error {
	var msg types.VehicleEventMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal vehicle message: %w", err)
	}

	releaseReq := types.ReleaseSlotRequest{
		SlotNumber: msg.SlotNumber,
		SlotType: types.GetSlotTypeByVehicleType(msg.VehicleType),
	}

	if err := s.integration.ReleaseSlot(ctx, msg.StructureUUID, msg.StructureType, releaseReq); err != nil {
		return fmt.Errorf("failed to release slot in structure: %w", err)
	}

	releaseLockReq := types.ReleaseSlotLockRequest{
		VehicleUUID: msg.VehicleUUID,
		StructureUUID: msg.StructureUUID,
		StructureSlotRequest: types.StructureSlotRequest{
			SlotNumber: msg.SlotNumber,
			SlotType: types.GetSlotTypeByVehicleType(msg.VehicleType),
		},
	}

	if err := s.integration.ReleaseSlotLock(ctx, releaseLockReq); err != nil {
		return fmt.Errorf("failed to release slot lock in tower leader: %w", err)
	}

	return nil
}
