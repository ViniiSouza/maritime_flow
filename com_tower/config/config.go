package config

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/utils"
)

type Config struct {
	Id uuid.UUID
	DB *pgx.Conn

	PropagationInterval time.Duration
	HeartbeatTimeout    time.Duration
}

func InitConfig(ctx context.Context) Config {
	id, err := uuid.Parse(os.Getenv(utils.TowerIdEnv))
	if err != nil {
		log.Fatalf("invalid id %s in env %s: %v", id, utils.TowerIdEnv, err)
	}

	conn, err := pgx.Connect(ctx, os.Getenv(utils.PostgresURIEnv))
	if err != nil {
		log.Fatalf("failed to establish database connection: %v", err)
	}

	interval, err := strconv.Atoi(os.Getenv(utils.PropagationIntervalEnv))
	if err != nil {
		log.Fatalf("failed to parse propagation interval env: %v", err)
	}

	propagationInterval := time.Duration(interval) * time.Second

	timeout, err := strconv.Atoi(os.Getenv(utils.HeartbeatTimeoutEnv))
	if err != nil {
		log.Fatalf("failed to parse heartbeat timeout env: %v", err)
	}

	heartbeatTimeout := time.Duration(timeout) * time.Second
	return Config{
		Id:                  id,
		DB:                  conn,
		PropagationInterval: propagationInterval,
		HeartbeatTimeout:    heartbeatTimeout,
	}
}
