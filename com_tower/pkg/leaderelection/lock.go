package leaderelection

import (
	"context"
	"log"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
)

const (
	AcquireLockQuery = "UPDATE tower_lock SET leader_id = $1, renewed_at = NOW() WHERE leader_id IS NULL;"
	GetLeaderQuery   = "SELECT leader_id FROM tower_lock LIMIT 1;"
)

func AcquireLockIfEmptyAndReturnLeaderUUID(ctx context.Context) types.UUID {
	tag, err := config.Configuration.GetDBConn().Exec(ctx, AcquireLockQuery, config.Configuration.GetId())
	if err != nil {
		log.Fatalf("failed to ensure leader lock: %v", err)
	}

	if tag.RowsAffected() == 0 {
		var id types.UUID
		if err := config.Configuration.GetDBConn().QueryRow(ctx, GetLeaderQuery).Scan(&id); err != nil {
			log.Fatalf("failed to query leader uuid: %v", err)
		}

		return id
	}

	return config.Configuration.GetId()
}
