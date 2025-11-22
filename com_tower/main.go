package main

import (
	"context"
	"log"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/leaderelection"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/tower/leader"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/tower/minion"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
)

var activeRoleCleanup func()

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config.InitConfig(ctx)
	leaderUuid := leaderelection.AcquireLockIfEmptyAndReturnLeaderUUID(ctx)
	config.Configuration.SetLeaderUUID(leaderUuid)

	if config.Configuration.IsLeader() {
		log.Println("starting initial role: LEADER")
		activeRoleCleanup = leader.InitLeader(ctx)
	} else {
		log.Println("starting initial role: MINION")
		activeRoleCleanup = minion.InitMinion(ctx) 
	}

main_loop:
	for {
		select {
		case role := <-leaderelection.ChangeRoleCh:
			log.Printf("received request to change role...")

			if activeRoleCleanup != nil {
				log.Println("stopping current active role goroutines...")
				activeRoleCleanup() 
			}

			if role == types.Leader {
				log.Printf("role requested: LEADER")
				activeRoleCleanup = leader.InitLeader(ctx)
				break
			}

			if role == types.Minion {
				log.Printf("role requested: MINION")
				activeRoleCleanup = minion.InitMinion(ctx)
			}
		case <-ctx.Done():
			break main_loop
		}
	}

	if activeRoleCleanup != nil {
		activeRoleCleanup()
	}
}
