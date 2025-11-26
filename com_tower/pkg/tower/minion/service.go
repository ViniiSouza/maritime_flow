package minion

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/utils"
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

func (s service) CheckSlotAvailability(ctx context.Context, request types.SlotRequest) (result *types.SlotResponse, err error) {
	failureCount := 0
	
	for (failureCount < config.Configuration.GetMaxStructureFailures()) {
		result, err = s.integration.RequestSlotToStructure(ctx, request)
		if err != nil {
			if errors.Is(err, utils.ErrStructureUnreachable) {
				failureCount++
				time.Sleep(500 * time.Millisecond)
				continue
			}
			return nil, err
		}

		break
	}

	if failureCount == config.Configuration.GetMaxStructureFailures() {
		log.Printf("failed to request slot to %s %s: %v", request.StructureType, request.StructureUUID.String(), err)
		if err := s.integration.SendEmail(request.StructureType, request.StructureUUID); err != nil {
			return nil, fmt.Errorf("failed to send email about structure failure: %w", err)
		}

		return &types.SlotResponse{
			State: types.InUseSlotState,
		}, nil
	}

	if result.State == types.FreeSlotState {
		acquireRequest := types.AcquireSlotRequest{
			VehicleUUID:          request.VehicleUUID,
			StructureUUID:        request.StructureUUID,
			StructureSlotRequest: request.StructureSlotRequest,
		}

		acquireResult, err := s.integration.AcquireSlotLockInTowerLeader(ctx, acquireRequest)
		if err != nil {
			log.Printf("failed to request slot to tower leader: %v", err)
			
			// rollback slot request in structure
			releaseSlotReq := types.ReleaseSlotRequest{SlotNumber: request.SlotNumber, SlotType: request.SlotType}
			if err := s.integration.ReleaseSlot(ctx, request.StructureUUID, request.StructureType, releaseSlotReq); err != nil {
				return nil, fmt.Errorf("failed to rollback slot request in %s %s: %w", request.StructureType, request.StructureUUID.String(), err)
			}
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
