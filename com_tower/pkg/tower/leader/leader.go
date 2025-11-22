package leader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/utils"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
)

var (
	client = &http.Client{}
)

func InitLeader(ctx context.Context) {
	go serve()
	go propagate(ctx)
}

func serve() {
	server := &http.Server{
		Handler:        setupRouter(),
		Addr:           fmt.Sprintf(":%s", os.Getenv(utils.PortEnv)),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func propagate(ctx context.Context) {
	repo := newRepository()
	service := newService(repo)

	for {
		select {
		case <-ctx.Done():
			return
		
		default:
			healthyTowers, err := service.ListHealthyTowers(ctx, config.Configuration.GetHeartbeatTimeout())
			if err != nil {
				log.Printf("[leader][propagate] failed to list healthy towers: %v", err)
				break
			}

			towersPayload, err := json.Marshal(types.TowersPayload{Towers: healthyTowers})
			if err != nil {
				log.Printf("[leader][propagate] failed to marshal healthy towers payload: %v", err)
				break
			}

			structures, err := service.ListStructures(ctx)
			if err != nil {
				log.Printf("[leader][propagate] failed to list structures: %v", err)
				break
			}

			structuresPayload, err := json.Marshal(structures)
			if err != nil {
				log.Printf("[leader][propagate] failed to marshal structures payload: %v", err)
				break
			}

			for _, tower := range healthyTowers {
				baseEndpoint := fmt.Sprintf("%s.tower.%s", tower.UUID.String(), config.Configuration.GetBaseDns())
				towersEndpoint := fmt.Sprintf("%s/%s", baseEndpoint, utils.TowersPropagationPath)
				structuresEndpoint := fmt.Sprintf("%s/%s", baseEndpoint, utils.StructuresPropagationPath)

				if err := doPropagateReq(ctx, towersEndpoint, towersPayload); err != nil {
					log.Printf("[leader][propagate] failed to propagate healthy towers to tower %s: %v", tower.UUID.String(), err)
				}

				if err := doPropagateReq(ctx, structuresEndpoint, structuresPayload); err != nil {
					log.Printf("[leader][propagate] failed to propagate structures to tower %s: %v", tower.UUID.String(), err)
				}
			}
		}

		time.Sleep(config.Configuration.GetPropagationInterval())
	}
}

func doPropagateReq(ctx  context.Context, endpoint string, payload []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create propagation request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute propagation request: %w", err)
	}

	defer resp.Body.Close()

	if _, err = io.Copy(io.Discard, resp.Body); err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	return nil
}
