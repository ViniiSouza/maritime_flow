package leader

import (
	"context"
	"fmt"
	"time"

	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/slot"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/structure"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/tower"
	"github.com/google/uuid"
)

type service struct {
	repository repository
}

func newService(r repository) service {
	return service{
		repository: r,
	}
}

func (s service) MarkTowerAsAlive(ctx context.Context, id uuid.UUID) (err error) {
	if _, err := s.repository.GetTowerById(ctx, id); err != nil {
		return fmt.Errorf("failed to check if tower exists: %w", err)
	}

	if err = s.repository.UpdateTowerLastSeen(ctx, id); err != nil {
		return fmt.Errorf("failed to mark tower as alive: %w", err)
	}

	return
}

func (s service) ListHealthyTowers(ctx context.Context, heartbeatTimeout time.Duration) (towers []tower.Tower, err error) {
	towers, err = s.repository.ListTowersByLastSeenAt(ctx, int(heartbeatTimeout.Seconds()))
	if err != nil {
		return nil, fmt.Errorf("failed to list towers: %w", err)
	}

	return
}

func (s service) ListStructures(ctx context.Context) (*structure.Structures, error) {
	platforms, err := s.repository.ListPlatforms(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list platforms: %w", err)
	}

	centrals, err := s.repository.ListCentrals(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list centrals: %w", err)
	}

	return &structure.Structures{
		Platforms: platforms,
		Centrals:  centrals,
	}, nil
}

func (s service) AcquireSlot(ctx context.Context, request slot.AcquireSlotRequest) (*slot.AcquireSlotResponse, error) {
	slotUuid, err := s.repository.GetSlotUUID(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get slot uuid: %w", err)
	}

	isSlotAvailable, err := s.repository.CheckSlotAvailability(ctx, slotUuid)
	if err != nil {
		return nil, fmt.Errorf("failed to check slot %s availability: %w", slotUuid, err)
	}

	if !isSlotAvailable {
		return &slot.AcquireSlotResponse{
			Result: slot.UnavailableAcquireSlotResultType,
		}, nil
	}

	if err := s.repository.AcquireSlot(ctx, request.VehicleUuid, slotUuid); err != nil {
		return nil, fmt.Errorf("failed to acquire slot %s: %w", slotUuid, err)
	}

	return &slot.AcquireSlotResponse{
		Result: slot.AcquiredAcquireSlotResultType,
	}, nil 
}
