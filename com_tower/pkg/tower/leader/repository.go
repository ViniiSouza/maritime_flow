package leader

import (
	"context"
	"errors"
	"strconv"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	listStructuresQuery = "SELECT st.id, st.latitude, st.longitude, jsonb_build_object('docks_qtt', COUNT(*) FILTER (WHERE sl.type = 'dock'), 'helipads_qtt', COUNT(*) FILTER (WHERE sl.type = 'helipad')) AS slots FROM structures st LEFT JOIN slots sl ON st.id = sl.structure_id WHERE st.type = $1 GROUP BY st.id;"
)

type repository struct {
	DB *pgxpool.Pool
}

func newRepository() repository {
	return repository{
		DB: config.Configuration.GetDBPool(),
	}
}

func (r repository) GetTowerById(ctx context.Context, id types.UUID) (types.Tower, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, latitude, longitude FROM towers WHERE id = $1;", id.String())
	if err != nil {
		return types.Tower{}, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[types.Tower])
}

func (r repository) UpdateTowerLastSeen(ctx context.Context, id types.UUID) (err error) {
	_, err = r.DB.Exec(ctx, "UPDATE towers SET last_seen_at = NOW() WHERE id = $1;", id.String())
	return
}

func (r repository) ListTowersByLastSeenAt(ctx context.Context, heartbeatTimeout int) ([]types.Tower, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, latitude, longitude FROM towers WHERE last_seen_at >= (NOW() - ($1 || ' seconds')::interval);", strconv.Itoa(heartbeatTimeout))
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[types.Tower])
}

func (r repository) ListPlatforms(ctx context.Context) ([]types.Platform, error) {
	rows, err := r.DB.Query(ctx, listStructuresQuery, "Platform")
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[types.Platform])
}

func (r repository) ListCentrals(ctx context.Context) ([]types.Central, error) {
	rows, err := r.DB.Query(ctx, listStructuresQuery, "Central")
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[types.Central])
}

func (r repository) GetSlotUUID(ctx context.Context, structureUuid types.UUID, slotType types.SlotType, slotNumber int) (slotUuid types.UUID, err error) {
	err = r.DB.QueryRow(ctx, "SELECT id FROM slots WHERE structure_id = $1 AND type = $2 AND number = $3;", structureUuid.String(), slotType, strconv.Itoa(slotNumber)).Scan(&slotUuid)
	return
}

func (r repository) CheckSlotAvailability(ctx context.Context, slotUuid types.UUID) (isAvailable bool, err error) {
	err = r.DB.QueryRow(ctx, "SELECT NOT EXISTS (SELECT 1 FROM vehicles WHERE current_slot_id = $1);", slotUuid.String()).Scan(&isAvailable)
	return
}

func (r repository) AcquireSlot(ctx context.Context, vehicleUuid types.UUID, slotUuid types.UUID) error {
	tag, err := r.DB.Exec(ctx, "UPDATE vehicles SET current_slot_id = $1 WHERE id = $2;", slotUuid.String(), vehicleUuid.String())
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("no rows affected, slot was not acquired")
	}

	return nil
}

func (r repository) ReleaseSlot(ctx context.Context, vehicleUuid types.UUID, slotUuid types.UUID) error {
	tag, err := r.DB.Exec(ctx, "UPDATE vehicles SET current_slot_id = NULL WHERE id = $1 AND current_slot_id = $2;", vehicleUuid.String(), slotUuid.String())
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("no rows affected, slot was not released")
	}

	return nil
}

func (r repository) AcquireLock(ctx context.Context) error {
	tag, err := r.DB.Exec(ctx, "UPDATE tower_lock SET leader_id = $1, renewed_at = NOW() WHERE leader_id = $1 OR leader_id IS NULL OR renewed_at < (NOW() - ($2 || ' seconds')::interval);", config.Configuration.GetIdAsString(), strconv.Itoa(int(config.Configuration.GetRenewLockTimeout().Seconds())))
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("no rows affected, lock was not acquired")
	}

	tag, err = r.DB.Exec(ctx, "UPDATE towers SET is_leader = (id = $1);", config.Configuration.GetIdAsString())
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("failed to set other towers as non leaders")
	}

	return nil
}

func (r repository) ReleaseLock(ctx context.Context) error {
	tag, err := r.DB.Exec(ctx, "UPDATE tower_lock SET leader_id = NULL WHERE leader_id = $1;", config.Configuration.GetIdAsString())
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("no rows affected, lock was not released")
	}

	return nil
}

func (r repository) RenewLock(ctx context.Context) error {
	tag, err := r.DB.Exec(ctx, "UPDATE tower_lock SET leader_id = $1, renewed_at = NOW() WHERE leader_id = $1;", config.Configuration.GetIdAsString())
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("no rows affected, lock was not renewed")
	}

	return nil
}
