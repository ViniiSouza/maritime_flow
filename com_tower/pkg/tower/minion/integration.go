package minion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/slot"
)

type integration struct {
	client *http.Client
	dns    string
}

func newIntegration() integration {
	return integration{
		client: &http.Client{},
		dns:    config.Configuration.GetBaseDns(),

	}
}

func (i integration) RequestSlotToStructure(ctx context.Context, slotRequest slot.SlotRequest) (slot.SlotResponse, error) {
	url := fmt.Sprintf("%s.%s.%s/slots", slotRequest.StructureUuid, slotRequest.StructureType, i.dns)
	payload, err := json.Marshal(slotRequest.StructureSlotRequest)
	if err != nil {
		return slot.SlotResponse{}, fmt.Errorf("failed to marshal slot request for %s %s: %w", slotRequest.StructureType, slotRequest.StructureUuid, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return slot.SlotResponse{}, fmt.Errorf("failed to create slot request for %s %s: %w", slotRequest.StructureType, slotRequest.StructureUuid, err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return slot.SlotResponse{}, fmt.Errorf("failed to request a slot for %s %s: %w", slotRequest.StructureType, slotRequest.StructureUuid, err)
	}

	defer resp.Body.Close()

	var slotResp slot.SlotResponse
	if err := json.NewDecoder(resp.Body).Decode(&slotResp); err != nil {
		return slot.SlotResponse{}, fmt.Errorf("failed to decode response body for  %s %s: %w", slotRequest.StructureType, slotRequest.StructureUuid, err)
	}

	return slotResp, nil
}

func (i integration) AcquireSlotLockInTowerLeader(ctx context.Context, slotRequest slot.SlotRequest) (slot.SlotResponse, error) {
	url := fmt.Sprintf("%s.tower.%s/acquire-slot", config.Configuration.GetLeaderUUID(), i.dns)
	payload, err := json.Marshal(slotRequest.StructureSlotRequest)
	if err != nil {
		return slot.SlotResponse{}, fmt.Errorf("failed to marshal slot request for %s %s: %w", slotRequest.StructureType, slotRequest.StructureUuid, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return slot.SlotResponse{}, fmt.Errorf("failed to create slot request for %s %s: %w", slotRequest.StructureType, slotRequest.StructureUuid, err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return slot.SlotResponse{}, fmt.Errorf("failed to request a slot for %s %s: %w", slotRequest.StructureType, slotRequest.StructureUuid, err)
	}

	defer resp.Body.Close()

	var slotResp slot.SlotResponse
	if err := json.NewDecoder(resp.Body).Decode(&slotResp); err != nil {
		return slot.SlotResponse{}, fmt.Errorf("failed to decode response body for  %s %s: %w", slotRequest.StructureType, slotRequest.StructureUuid, err)
	}

	return slotResp, nil
}
