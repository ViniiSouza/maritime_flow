package leaderelection

import (
	"context"
	"log"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
	"github.com/google/uuid"
)

const (
	AcquireLockQuery = "UPDATE tower_lock SET leader_id = $1, renewed_at = NOW() WHERE leader_id IS NULL;"
	GetLeaderQuery   = "SELECT leader_id FROM tower_lock LIMIT 1;"
)

func AcquireLockIfEmptyAndReturnLeaderUUID(ctx context.Context) types.UUID {
	tag, err := config.Configuration.GetDBPool().Exec(ctx, AcquireLockQuery, config.Configuration.GetId())
	if err != nil {
		log.Fatalf("failed to ensure leader lock: %v", err)
	}

	if tag.RowsAffected() == 0 {
		var id string
		if err := config.Configuration.GetDBPool().QueryRow(ctx, GetLeaderQuery).Scan(&id); err != nil {
			log.Fatalf("failed to query leader uuid: %v", err)
		}

		uuid, err := uuid.Parse(id)
		if err != nil {
			log.Fatalf("failed to parse leader uuid: %v", err)
		}

		return types.UUID(uuid)
	}

	return config.Configuration.GetId()
}
