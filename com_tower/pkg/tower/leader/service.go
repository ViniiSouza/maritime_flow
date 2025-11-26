package leader

import (
	"context"
	"fmt"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
)

type service struct {
	repository repository
}

func newService(r repository) service {
	return service{
		repository: r,
	}
}

func (s service) AcquireLock(ctx context.Context) error {
	return s.repository.AcquireLock(ctx)
}

func (s service) ReleaseLock(ctx context.Context) error {
	return s.repository.ReleaseLock(ctx)
}

func (s service) RenewLock(ctx context.Context) error {
	return s.repository.RenewLock(ctx)
}

func (s service) MarkTowerAsAlive(ctx context.Context, id types.UUID) (err error) {
	if _, err := s.repository.GetTowerById(ctx, id); err != nil {
		return fmt.Errorf("failed to check if tower exists: %w", err)
	}

	if err = s.repository.UpdateTowerLastSeen(ctx, id); err != nil {
		return fmt.Errorf("failed to mark tower as alive: %w", err)
	}

	return
}

func (s service) ListHealthyTowers(ctx context.Context) (towers []types.Tower, err error) {
	towers, err = s.repository.ListTowersByLastSeenAt(ctx, int(config.Configuration.GetHeartbeatTimeout().Seconds()))
	if err != nil {
		return nil, err
	}

	return
}

func (s service) ListStructures(ctx context.Context) (*types.Structures, error) {
	platforms, err := s.repository.ListPlatforms(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list platforms: %w", err)
	}

	centrals, err := s.repository.ListCentrals(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list centrals: %w", err)
	}

	return &types.Structures{
		Platforms: platforms,
		Centrals:  centrals,
	}, nil
}

func (s service) AcquireSlot(ctx context.Context, request types.AcquireSlotRequest) (*types.AcquireSlotResponse, error) {
	slotUuid, err := s.repository.GetSlotUUID(ctx, request.StructureUUID, request.SlotType, request.SlotNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get slot uuid: %w", err)
	}

	isSlotAvailable, err := s.repository.CheckSlotAvailability(ctx, slotUuid)
	if err != nil {
		return nil, fmt.Errorf("failed to check slot %s availability: %w", slotUuid.String(), err)
	}

	if !isSlotAvailable {
		return &types.AcquireSlotResponse{
			Result: types.UnavailableAcquireSlotResultType,
		}, nil
	}

	if err := s.repository.AcquireSlot(ctx, request.VehicleUUID, slotUuid); err != nil {
		return nil, fmt.Errorf("failed to acquire slot %s: %w", slotUuid.String(), err)
	}

	return &types.AcquireSlotResponse{
		Result: types.AcquiredAcquireSlotResultType,
	}, nil 
}

func (s service) ReleaseSlot(ctx context.Context, request types.ReleaseSlotLockRequest) error {
	slotUuid, err := s.repository.GetSlotUUID(ctx, request.StructureUUID, request.SlotType, request.SlotNumber)
	if err != nil {
		return fmt.Errorf("failed to get slot uuid: %w", err)
	}

	if err := s.repository.ReleaseSlot(ctx, request.VehicleUUID, slotUuid); err != nil {
		return fmt.Errorf("failed to release slot %s: %w", slotUuid.String(), err)
	}

	return nil 
}
