package leader

import (
	"context"
	"errors"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
	"github.com/jackc/pgx/v5"
)

type repository struct {
	DB *pgx.Conn
}

func newRepository() repository {
	return repository{
		DB: config.Configuration.GetDBConn(),
	}
}

func (r repository) GetTowerById(ctx context.Context, id types.UUID) (types.Tower, error) {
	rows, err := r.DB.Query(ctx, "SELECT id FROM towers WHERE id == $1;", id.String())
	if err != nil {
		return types.Tower{}, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[types.Tower])
}

func (r repository) UpdateTowerLastSeen(ctx context.Context, id types.UUID) (err error) {
	_, err = r.DB.Exec(ctx, "UPDATE towers SET last_seen = NOW() WHERE id == $1;", id.String())
	return
}

func (r repository) ListTowersByLastSeenAt(ctx context.Context, heartbeatTimeout int) ([]types.Tower, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, latitude, longitude FROM towers WHERE last_seen_at >= (NOW() - ($1 || ' seconds')::interval);", heartbeatTimeout)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[types.Tower])
}

func (r repository) ListPlatforms(ctx context.Context) ([]types.Platform, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, latitude, longitude FROM platforms;")
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[types.Platform])
}

func (r repository) ListCentrals(ctx context.Context) ([]types.Central, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, latitude, longitude FROM centrals;")
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[types.Central])
}

func (r repository) GetSlotUUID(ctx context.Context, structureUuid types.UUID, slotType types.SlotType, slotNumber int) (slotUuid types.UUID, err error) {
	err = r.DB.QueryRow(ctx, "SELECT id FROM slots WHERE structure_id = $1 AND type = $2 AND number = $3;", structureUuid, slotType, slotNumber).Scan(&slotUuid)
	return
}

func (r repository) CheckSlotAvailability(ctx context.Context, slotUuid types.UUID) (isAvailable bool, err error) {
	err = r.DB.QueryRow(ctx, "SELECT NOT EXISTS (SELECT 1 FROM vehicles WHERE current_slot_uuid = $1);", slotUuid).Scan(&isAvailable)
	return
}

func (r repository) AcquireSlot(ctx context.Context, vehicleUuid types.UUID, slotUuid types.UUID) error {
	tag, err := r.DB.Exec(ctx, "UPDATE vehicles SET current_slot_id = $1 WHERE id = $2;", slotUuid, vehicleUuid)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("no rows affected, slot was not acquired")
	}

	return nil
}

func (r repository) ReleaseSlot(ctx context.Context, vehicleUuid types.UUID, slotUuid types.UUID) error {
	tag, err := r.DB.Exec(ctx, "UPDATE vehicles SET current_slot_id = NULL WHERE id = $1 AND current_slot_id = $2;", vehicleUuid, slotUuid)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("no rows affected, slot was not released")
	}

	return nil
}
