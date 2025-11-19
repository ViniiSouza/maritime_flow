package leader

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/utils"
)

func InitLeader(ctx context.Context, cfg config.Config) error {
	go serve(cfg)

	go propagate(ctx, cfg)

	return nil
}

func serve(cfg config.Config) {
	server := &http.Server{
		Handler:        setupRouter(cfg),
		Addr:           fmt.Sprintf(":%s", os.Getenv(utils.PortEnv)),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func propagate(ctx context.Context, cfg config.Config) {
	repo := newRepository(cfg.DB)
	service := newService(repo)

	for {
		select {
		case <-ctx.Done():
			return
		
		default:
			healthyTowers, err := service.ListHealthyTowers(ctx, cfg.HeartbeatTimeout)
			if err != nil {
				log.Printf("[leader][propagate] failed to list healthy towers: %v", err)
				break
			}

			structures, err := service.ListStructures(ctx)
			if err != nil {
				log.Printf("[leader][propagate] failed to list structures: %v", err)
				break
			}
		}

		time.Sleep(cfg.PropagationInterval)
	}
} 