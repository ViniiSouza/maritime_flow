package leader

import (
	"context"

	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/structure"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/tower"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type repository struct {
	DB *pgx.Conn
}

func newRepository(db *pgx.Conn) repository {
	return repository{
		DB: db,
	}
}

func (r repository) GetTowerById(ctx context.Context, id uuid.UUID) (tower.Tower, error) {
	rows, err := r.DB.Query(ctx, "SELECT id FROM towers WHERE id == $1", id.String())
	if err != nil {
		return tower.Tower{}, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[tower.Tower])
}

func (r repository) UpdateTowerLastSeen(ctx context.Context, id uuid.UUID) (err error) {
	_, err = r.DB.Exec(ctx, "UPDATE towers SET last_seen = NOW() WHERE id == $1", id.String())
	return
}

func (r repository) ListTowersByLastSeenAt(ctx context.Context, heartbeatTimeout int) ([]tower.Tower, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, latitude, longitude FROM towers WHERE last_seen_at >= NOW() - ($1 || ' seconds')::interval);", heartbeatTimeout)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[tower.Tower])
}

func (r repository) ListPlatforms(ctx context.Context) ([]structure.Platform, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, latitude, longitude FROM platforms")
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[structure.Platform])
}

func (r repository) ListCentrals(ctx context.Context) ([]structure.Central, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, latitude, longitude FROM centrals")
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[structure.Central])
}
