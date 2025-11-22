package minion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/slot"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/structure"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/tower"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/utils"
)

type integration struct {
	client *http.Client
}

func newIntegration() integration {
	return integration{
		client: &http.Client{},
	}
}

func (i integration) RequestSlotToStructure(ctx context.Context, slotRequest slot.SlotRequest) (*slot.SlotResponse, error) {
	url := fmt.Sprintf("%s.%s.%s/slots", slotRequest.StructureUUID, slotRequest.StructureType, config.Configuration.GetBaseDns())
	payload, err := json.Marshal(slotRequest.StructureSlotRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal slot request for %s %s: %w", slotRequest.StructureType, slotRequest.StructureUUID, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create slot request for %s %s: %w", slotRequest.StructureType, slotRequest.StructureUUID, err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request a slot for %s %s: %w", slotRequest.StructureType, slotRequest.StructureUUID, err)
	}

	defer resp.Body.Close()

	var slotResp slot.SlotResponse
	if err := json.NewDecoder(resp.Body).Decode(&slotResp); err != nil {
		return nil, fmt.Errorf("failed to decode response body for %s %s: %w", slotRequest.StructureType, slotRequest.StructureUUID, err)
	}

	return &slotResp, nil
}

func (i integration) AcquireSlotLockInTowerLeader(ctx context.Context, slotRequest slot.AcquireSlotRequest) (*slot.AcquireSlotResponse, error) {
	url := fmt.Sprintf("%s.tower.%s/acquire-slot", config.Configuration.GetLeaderUUID(), config.Configuration.GetBaseDns())
	payload, err := json.Marshal(slotRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal slot acquire request for %s %d in structure %s: %w", slotRequest.SlotType, slotRequest.SlotNumber, slotRequest.StructureUUID, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create slot acquire request for %s %d in structure %s: %w", slotRequest.SlotType, slotRequest.SlotNumber, slotRequest.StructureUUID, err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request a slot acquire for %s %d in structure %s: %w", slotRequest.SlotType, slotRequest.SlotNumber, slotRequest.StructureUUID, err)
	}

	defer resp.Body.Close()

	var acquireResp slot.AcquireSlotResponse
	if err := json.NewDecoder(resp.Body).Decode(&acquireResp); err != nil {
		return nil, fmt.Errorf("failed to decode slot acquire response body for %s %d in structure %s: %w", slotRequest.SlotType, slotRequest.SlotNumber, slotRequest.StructureUUID, err)
	}

	return &acquireResp, nil
}

func (i integration) SendHealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s.tower.%s/tower-health", config.Configuration.GetLeaderUUID(), config.Configuration.GetBaseDns())
	payload, err := json.Marshal(tower.TowerHealthRequest{Id: config.Configuration.GetId()})
	if err != nil {
		return fmt.Errorf("failed to marshal healthcheck request for tower %s: %w", config.Configuration.GetIdAsString(), err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create healthcheck request for tower %s: %w", config.Configuration.GetIdAsString(), err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute healthcheck request for tower %s: %w", config.Configuration.GetIdAsString(), err)
	}

	defer resp.Body.Close()

	if _, err = io.Copy(io.Discard, resp.Body); err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	return nil
}

func (i integration) ReleaseSlot(ctx context.Context, structureUuid utils.UUID, structureType structure.StructureType, slotRequest slot.ReleaseSlotRequest) error {
	url := fmt.Sprintf("%s.%s.%s/release-slot", structureUuid, structureType, config.Configuration.GetBaseDns())
	payload, err := json.Marshal(slotRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal release slot request for %s %s: %w", structureType, structureUuid, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create release slot request for %s %s: %w", structureType, structureUuid, err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to request a slot release for %s %s: %w", structureType, structureUuid, err)
	}

	defer resp.Body.Close()

	if _, err = io.Copy(io.Discard, resp.Body); err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	return nil
}

func (i integration) ReleaseSlotLock(ctx context.Context, slotRequest slot.ReleaseSlotLockRequest) error {
	url := fmt.Sprintf("%s.tower.%s/release-slot", config.Configuration.GetLeaderUUID(), config.Configuration.GetBaseDns())
	payload, err := json.Marshal(slotRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal release slot request for tower %s: %w", config.Configuration.GetIdAsString(), err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create release slot request for tower %s: %w", config.Configuration.GetIdAsString(), err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute release slot request for tower %s: %w", config.Configuration.GetIdAsString(), err)
	}

	defer resp.Body.Close()

	if _, err = io.Copy(io.Discard, resp.Body); err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	return nil
}
