package minion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
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

func (i integration) RequestSlotToStructure(ctx context.Context, slotRequest types.SlotRequest) (*types.SlotResponse, error) {
	url := fmt.Sprintf("http://s-%s.%s.%s/slots", slotRequest.StructureUUID.String(), slotRequest.StructureType, config.Configuration.GetBaseDns())
	payload, err := json.Marshal(slotRequest.StructureSlotRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal slot request for %s %s: %w", slotRequest.StructureType, slotRequest.StructureUUID.String(), err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create slot request for %s %s: %w", slotRequest.StructureType, slotRequest.StructureUUID.String(), err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request a slot for %s %s: %w: %w", slotRequest.StructureType, slotRequest.StructureUUID.String(), utils.ErrStructureUnreachable, err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var slotResp types.SlotResponse
		if err := json.NewDecoder(resp.Body).Decode(&slotResp); err != nil {
			return nil, fmt.Errorf("failed to decode response body for %s %s: %w", slotRequest.StructureType, slotRequest.StructureUUID.String(), err)
		}

		return &slotResp, nil

	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return nil, fmt.Errorf("failed to request a slot for %s %s: %w: %w", slotRequest.StructureType, slotRequest.StructureUUID.String(), utils.ErrStructureUnreachable, err)

	default:
		return nil, utils.HttpErrorNotHandled(resp.StatusCode, resp.Body)
	}
}

func (i integration) AcquireSlotLockInTowerLeader(ctx context.Context, slotRequest types.AcquireSlotRequest) (*types.AcquireSlotResponse, error) {
	url := fmt.Sprintf("http://t-%s.tower.%s/acquire-slot", config.Configuration.GetLeaderUUIDAsString(), config.Configuration.GetBaseDns())
	payload, err := json.Marshal(slotRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal slot acquire request for %s %d in structure %s: %w", slotRequest.SlotType, slotRequest.SlotNumber, slotRequest.StructureUUID.String(), err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create slot acquire request for %s %d in structure %s: %w", slotRequest.SlotType, slotRequest.SlotNumber, slotRequest.StructureUUID.String(), err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request a slot acquire for %s %d in structure %s: %w", slotRequest.SlotType, slotRequest.SlotNumber, slotRequest.StructureUUID.String(), err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var acquireResp types.AcquireSlotResponse
		if err := json.NewDecoder(resp.Body).Decode(&acquireResp); err != nil {
			return nil, fmt.Errorf("failed to decode slot acquire response body for %s %d in structure %s: %w", slotRequest.SlotType, slotRequest.SlotNumber, slotRequest.StructureUUID.String(), err)
		}

		return &acquireResp, nil

	default:
		return nil, utils.HttpErrorNotHandled(resp.StatusCode, resp.Body)
	}
}

func (i integration) SendHealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("http://t-%s.tower.%s/tower-health", config.Configuration.GetLeaderUUIDAsString(), config.Configuration.GetBaseDns())
	payload, err := json.Marshal(types.TowerHealthRequest{Id: config.Configuration.GetId()})
	if err != nil {
		return fmt.Errorf("failed to marshal healthcheck request for tower %s: %w", config.Configuration.GetIdAsString(), err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create healthcheck request for tower %s: %w", config.Configuration.GetIdAsString(), err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return fmt.Errorf("tower %s: %w: %w", config.Configuration.GetIdAsString(), utils.ErrLeaderUnreachable, err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNoContent:
		if _, err = io.Copy(io.Discard, resp.Body); err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		return nil

	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return fmt.Errorf("tower %s: %w: http error code %d", config.Configuration.GetIdAsString(), utils.ErrLeaderUnreachable, resp.StatusCode) 

	default:
		return utils.HttpErrorNotHandled(resp.StatusCode, resp.Body)
	}
}

func (i integration) ReleaseSlot(ctx context.Context, structureUuid types.UUID, structureType types.StructureType, slotRequest types.ReleaseSlotRequest) error {
	url := fmt.Sprintf("http://s-%s.%s.%s/release-slot", structureUuid.String(), structureType, config.Configuration.GetBaseDns())
	payload, err := json.Marshal(slotRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal release slot request for %s %s: %w", structureType, structureUuid.String(), err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create release slot request for %s %s: %w", structureType, structureUuid.String(), err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to request a slot release for %s %s: %w", structureType, structureUuid.String(), err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNoContent:
		if _, err = io.Copy(io.Discard, resp.Body); err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		return nil

	default:
		return utils.HttpErrorNotHandled(resp.StatusCode, resp.Body)
	}
}

func (i integration) ReleaseSlotLock(ctx context.Context, slotRequest types.ReleaseSlotLockRequest) error {
	url := fmt.Sprintf("http://t-%s.tower.%s/release-slot", config.Configuration.GetLeaderUUIDAsString(), config.Configuration.GetBaseDns())
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

	switch resp.StatusCode {
	case http.StatusNoContent:
		if _, err = io.Copy(io.Discard, resp.Body); err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		return nil

	default:
		return utils.HttpErrorNotHandled(resp.StatusCode, resp.Body)
	}
}

func (i integration) SendEmail(structureType types.StructureType, structureUuid types.UUID) error {
	emailConfig := config.Configuration.GetEmailConfig()
	auth := smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.Host)

	subject := fmt.Sprintf(utils.EmailSubjectTemplate, string(structureType), structureUuid.String())
	message := fmt.Sprintf(utils.EmailBodyTemplate, config.Configuration.GetIdAsString(), string(structureType), structureUuid.String())
	toHeader := strings.Join(emailConfig.Recipients, ",")
	msgString := fmt.Sprintf("From: %s\r\n", emailConfig.Username) +
		fmt.Sprintf("To: %s\r\n", toHeader) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
		"\r\n" + // Empty line separates headers from body
		message

	address := emailConfig.Host + ":" + emailConfig.Port

	err := smtp.SendMail(address, auth, emailConfig.Username, emailConfig.Recipients, []byte(msgString))
	if err != nil {
		return err
	}

	return nil
}

