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
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/utils"
)

var (
	client = &http.Client{}
)

func InitLeader(ctx context.Context) func() {
	leaderCtx, leaderCancel := context.WithCancel(ctx)

	repo := newRepository()
	svc := newService(repo)
	if err := svc.AcquireLock(leaderCtx); err != nil {
		log.Fatalf("[leader] failed to acquire database lock: %v", err)
	}

	server := &http.Server{
		Handler:        setupRouter(svc),
		Addr:           fmt.Sprintf(":%s", os.Getenv(utils.PortEnv)),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go serve(server)
	go propagate(leaderCtx, svc)
	go renewLock(leaderCtx, svc)

	return func() {
		releaseCtx, releaseCancel := context.WithTimeout(context.Background(), 2*time.Second)
    	defer releaseCancel()

		if err := svc.ReleaseLock(releaseCtx); err != nil {
			log.Printf("[leader] failed to release database lock: %v", err)
		}
		leaderCancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("[leader] HTTP Server forced shutdown: %v", err)
		} else {
			log.Println("[leader] HTTP Server stopped")
		}
	}
}

func serve(server *http.Server) {
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func propagate(ctx context.Context, svc service) {
	for {
		select {
		case <-time.After(config.Configuration.GetPropagationInterval()):
			healthyTowers, err := svc.ListHealthyTowers(ctx, config.Configuration.GetHeartbeatTimeout())
			if err != nil {
				log.Printf("[leader][propagate] failed to list healthy towers: %v", err)
				break
			}

			towersPayload, err := json.Marshal(types.TowersPayload{Towers: healthyTowers})
			if err != nil {
				log.Printf("[leader][propagate] failed to marshal healthy towers payload: %v", err)
				break
			}

			structures, err := svc.ListStructures(ctx)
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
				baseEndpoint := fmt.Sprintf("t-%s.tower.%s", tower.UUID.String(), config.Configuration.GetBaseDns())
				towersEndpoint := fmt.Sprintf("%s/%s", baseEndpoint, utils.TowersPropagationPath)
				structuresEndpoint := fmt.Sprintf("%s/%s", baseEndpoint, utils.StructuresPropagationPath)

				if err := doPropagateReq(ctx, towersEndpoint, towersPayload); err != nil {
					log.Printf("[leader][propagate] failed to propagate healthy towers to tower %s: %v", tower.UUID.String(), err)
				}

				if err := doPropagateReq(ctx, structuresEndpoint, structuresPayload); err != nil {
					log.Printf("[leader][propagate] failed to propagate structures to tower %s: %v", tower.UUID.String(), err)
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

func renewLock(ctx context.Context, svc service) {
	for {
		select {
		case <-time.After(config.Configuration.GetRenewLockInterval()):
			if err := svc.RenewLock(ctx); err != nil {
				log.Printf("[leader][renew_lock] failed to renew lock: %v", err)
			}

		case <-ctx.Done():
			return
		}
	}
}

func doPropagateReq(ctx context.Context, endpoint string, payload []byte) error {
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
