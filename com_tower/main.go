package main

import (
	"context"
	"log"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/leaderelection"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/tower/leader"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/tower/minion"
)

func main() {
	ctx := context.Background()
	cfg := config.InitConfig(ctx)
	isLeader := leaderelection.AcquireLock(ctx, cfg)

	if isLeader {
		if err := leader.InitLeader(ctx, cfg); err != nil {
			log.Fatalf("failed to initialize leader tower: %v", err)
		}
	} else {
		if err := minion.InitMinion(); err != nil {
			log.Fatalf("failed to initialize minion tower: %v", err)
		}
	}

	<-ctx.Done()
}
