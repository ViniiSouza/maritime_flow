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
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/tower"
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
	url := fmt.Sprintf("%s.%s.%s/slots", slotRequest.StructureUuid, slotRequest.StructureType, config.Configuration.GetBaseDns())
	payload, err := json.Marshal(slotRequest.StructureSlotRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal slot request for %s %s: %w", slotRequest.StructureType, slotRequest.StructureUuid, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create slot request for %s %s: %w", slotRequest.StructureType, slotRequest.StructureUuid, err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request a slot for %s %s: %w", slotRequest.StructureType, slotRequest.StructureUuid, err)
	}

	defer resp.Body.Close()

	var slotResp slot.SlotResponse
	if err := json.NewDecoder(resp.Body).Decode(&slotResp); err != nil {
		return nil, fmt.Errorf("failed to decode response body for %s %s: %w", slotRequest.StructureType, slotRequest.StructureUuid, err)
	}

	return &slotResp, nil
}

func (i integration) AcquireSlotLockInTowerLeader(ctx context.Context, slotRequest slot.AcquireSlotRequest) (*slot.AcquireSlotResponse, error) {
	url := fmt.Sprintf("%s.tower.%s/acquire-slot", config.Configuration.GetLeaderUUID(), config.Configuration.GetBaseDns())
	payload, err := json.Marshal(slotRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal slot acquire request for %s %d in structure %s: %w", slotRequest.SlotType, slotRequest.SlotNumber, slotRequest.StructureUuid, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create slot acquire request for %s %d in structure %s: %w", slotRequest.SlotType, slotRequest.SlotNumber, slotRequest.StructureUuid, err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request a slot acquire for %s %d in structure %s: %w", slotRequest.SlotType, slotRequest.SlotNumber, slotRequest.StructureUuid, err)
	}

	defer resp.Body.Close()

	var acquireResp slot.AcquireSlotResponse
	if err := json.NewDecoder(resp.Body).Decode(&acquireResp); err != nil {
		return nil, fmt.Errorf("failed to decode slot acquire response body for %s %d in structure %s: %w", slotRequest.SlotType, slotRequest.SlotNumber, slotRequest.StructureUuid, err)
	}

	return &acquireResp, nil
}

func (i integration) SendHealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s.tower.%s/tower-health", config.Configuration.GetLeaderUUID(), config.Configuration.GetBaseDns())
	payload, err := json.Marshal(tower.TowerHealthRequest{Id: config.Configuration.GetIdAsString()})
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
