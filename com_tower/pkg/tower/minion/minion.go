package minion

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/leaderelection"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/utils"
)

func InitMinion(ctx context.Context) func() {
	minionCtx, minionCancel := context.WithCancel(ctx)

	integ := newIntegration()
	repo := newRepository()
	svc := newService(integ, repo)

	if err := bindAuditQueue(); err != nil {
		log.Fatalf("[minion][audit] failed to bind audit queue: %v", err)
	}
	
	server := &http.Server{
		Handler:        setupRouter(svc),
		Addr:           fmt.Sprintf(":%s", os.Getenv(utils.PortEnv)),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go serve(server)
	go healthcheck(minionCtx, svc)
	go consumeBroker(minionCtx, svc)

	return func() {
		minionCancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("[minion] HTTP Server forced shutdown: %v", err)
		} else {
			log.Println("[minion] HTTP Server stopped")
		}
	}
}

func serve(server *http.Server) {
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func healthcheck(ctx context.Context, svc service) {
	maxLeaderFailures := config.Configuration.GetMaxLeaderFailures()
	failureCount := 0

	for {
		select {
		case <-time.After(config.Configuration.GetHeartbeatInterval()):
			if err := svc.SendHealthCheck(ctx); err != nil {
				log.Printf("[minion][healthcheck] failed to send healthcheck: %v", err)

				if errors.Is(err, utils.ErrLeaderUnreachable) {
					failureCount++

					if failureCount == maxLeaderFailures {
						go leaderelection.StartElection(svc.ListTowers())
						failureCount = 0
						time.Sleep(5 * time.Second)
					}
				}
			} else {
				failureCount = 0
			}

		case <-ctx.Done():
			return
		}
	}
}

func consumeBroker(ctx context.Context, svc service) {
	slotReleaseCh, err := bindTowersQueue()
	if err != nil {
		log.Fatalf("[minion][consumer] failed to bind towers queue: %v", err)
	}

	for {
		select {
		case msg := <-slotReleaseCh:
			log.Printf("[minion][consumer] received message: %s", string(msg.Body))
			if err := svc.ReleaseSlot(ctx, msg.Body); err != nil {
				log.Printf("[minion][consumer] failed to release slot: %v", err)
			}

		case <-ctx.Done():
			log.Printf("[minion][consumer] interrupting consumer...")
			config.CloseRabbitMQ()
			return
		}
	}
}
