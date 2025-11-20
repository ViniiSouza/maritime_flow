package minion

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/utils"
)

func InitMinion(ctx context.Context) error {
	go serve()

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
