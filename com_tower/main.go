package main

import (
	"context"

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
		leader.InitLeader(ctx)
	} else {
		minion.InitMinion(ctx)
	}

	<-ctx.Done()
}
