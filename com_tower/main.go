package main

import (
	"context"
	"log"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/leaderelection"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/tower/leader"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/tower/minion"
)

var (
	isLeader = false
)

func main() {
	ctx := context.Background()
	config.InitConfig(ctx)
	leaderUuid := leaderelection.TryAcquireLockAndReturnLeaderUUID(ctx)
	config.Configuration.SetLeaderUUID(leaderUuid)
	if leaderUuid == config.Configuration.GetId() {
		isLeader = true
	} 

	if isLeader {
		if err := leader.InitLeader(ctx); err != nil {
			log.Fatalf("failed to initialize leader tower: %v", err)
		}
	} else {
		if err := minion.InitMinion(ctx); err != nil {
			log.Fatalf("failed to initialize minion tower: %v", err)
		}
	}

	<-ctx.Done()
}
