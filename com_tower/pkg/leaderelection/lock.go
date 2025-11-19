package leaderelection

import (
	"context"
	"log"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
)

const (
	AcquireLockQuery = "UPDATE tower_lock SET leader_id = ? WHERE leader_id IS NULL;"
)

func AcquireLock(ctx context.Context ,cfg config.Config) bool {
	tag, err := cfg.DB.Exec(ctx, AcquireLockQuery, cfg.Id)
	if err != nil {
		log.Fatalf("failed to ensure leader lock: %v", err)
	}

	if tag.RowsAffected() == 0 {
		return false
	}

	return true
}
