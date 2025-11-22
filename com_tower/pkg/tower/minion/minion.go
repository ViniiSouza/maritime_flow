package minion

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

func InitMinion(ctx context.Context) error {
	go serve()
	go healthcheck(ctx)

	return nil
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

func healthcheck(ctx context.Context) {
	integ := newIntegration()
	repo := newRepository()
	svc := newService(integ, repo)

	for {
		if err := svc.SendHealthCheck(ctx); err != nil {
			log.Printf("[minion][healthcheck] failed to send healthcheck: %v", err)
		}

		time.Sleep(config.Configuration.GetHeartbeatInterval())
	}
}
