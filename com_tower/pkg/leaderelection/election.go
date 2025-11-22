package leaderelection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
)

var ChangeRoleCh = make(chan types.Role)

func StartElection(towers []types.Tower) {
    uptime := config.Configuration.GetUptimeSeconds()
    log.Printf("[minion][election] starting leader election: my uptime: %.2fs", uptime)

    electionReq := types.ElectionRequest{
        CandidateUptime: uptime,
    }

    hasHighestUptime := true
    for _, tower := range towers {
        if config.Configuration.GetId() == tower.UUID {
            continue
        }

        url := fmt.Sprintf("%s.tower.%s/election", tower.UUID.String(), config.Configuration.GetBaseDns())
        payload, err := json.Marshal(electionReq)
		if err != nil {
			log.Printf("[minion][election] failed to marshal election request: %v", err)
			return
		}
        
        resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
        if err != nil {
            log.Printf("[minion][election] failed to send election request to tower %s: %v", tower.UUID, err)
            continue
        }

        defer resp.Body.Close()

        var electionResp types.ElectionResponse
        if err := json.NewDecoder(resp.Body).Decode(&electionResp); err != nil {
            log.Printf("[minion][election] failed to decode response from tower %s: %v", tower.UUID, err)
            continue
        }

        if resp.StatusCode == http.StatusOK && electionResp.HasHigherUptime {
            log.Printf("[minion][election] tower %s has a higher uptime of (%.2fs): stopping election", tower.UUID, electionResp.Uptime)
            hasHighestUptime = false
            break
        }
    }
    
    if !hasHighestUptime {
		log.Printf("[minion][election] election lost, delegating election to another tower with a higher uptime")
		return
    }
    
	log.Printf("[minion][election] election won, becoming leader")
	config.Configuration.SetLeaderUUID(config.Configuration.GetId())
	ChangeRoleCh <- types.Leader
	go broadcastCoordinator(towers)
}

func broadcastCoordinator(towers []types.Tower) {
    coordinatorReq := types.NewLeaderRequest{
        NewLeaderUUID: config.Configuration.GetId(),
    }

    payload, err := json.Marshal(coordinatorReq)
	if err != nil {
		log.Printf("[leader][election] failed to marshal new leader request: %v", err)
		return
	}

    for _, tower := range towers {
        if config.Configuration.GetId() == tower.UUID {
            continue
        }

        url := fmt.Sprintf("%s.tower.%s/leader", tower.UUID.String(), config.Configuration.GetBaseDns())
        resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
        if err != nil {
            log.Printf("[leader][election] failed to announce new leader to tower %s: %v", tower.UUID, err)
            continue
        }

        resp.Body.Close()
        log.Printf("[leader][election] announced new leader to tower %s.", tower.UUID)
    }
}
