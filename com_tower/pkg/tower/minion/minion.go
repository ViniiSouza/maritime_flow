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

func InitMinion(ctx context.Context) {
	go serve()
	go healthcheck(ctx)
	go consumeBroker(ctx)
}

func serve() {
	if err := bindAuditQueue(); err != nil {
		log.Fatalf("[minion][audit] failed to bind audit queue: %v", err)
	}

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

func consumeBroker(ctx context.Context) {
	integ := newIntegration()
	repo := newRepository()
	svc := newService(integ, repo)

	slotReleaseCh, err := bindTowersQueue()
	if err != nil {
		log.Fatalf("[minion][consumer] failed to bind towers queue: %v", err)
	}

main_loop:
	for {
		select {
		case msg := <-slotReleaseCh:
			log.Printf("received message: %s", string(msg.Body))
			svc.ReleaseSlot(ctx, msg.Body)

		case <-ctx.Done():
			log.Printf("interrupting...")
			config.CloseRabbitMQ()
			break main_loop
		}
	}
}
