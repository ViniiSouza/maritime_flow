package config

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/utils"
)

var brokerconn *amqp.Connection
var Configuration *Config

type Config struct {
	id          types.UUID
	leaderUuid  types.UUID
	baseDns     string
	towersQueue string
	auditQueue  string

	db       *pgxpool.Pool
	rabbitmq *amqp.Channel

	uptime time.Time

	maxLeaderFailures   int
	propagationInterval time.Duration
	heartbeatInterval   time.Duration
	heartbeatTimeout    time.Duration
	renewLockInterval   time.Duration
	renewLockTimeout    time.Duration
}

func (c *Config) GetId() types.UUID {
	return c.id
}

func (c *Config) GetIdAsString() string {
	return c.id.String()
}

func (c *Config) GetDBPool() *pgxpool.Pool {
	return c.db
}

func (c *Config) GetRabbitMQChannel() *amqp.Channel {
	return c.rabbitmq
}

func (c *Config) GetTowersQueue() string {
	return c.towersQueue
}

func (c *Config) GetAuditQueue() string {
	return c.auditQueue
}

func (c *Config) GetBaseDns() string {
	return c.baseDns
}

func (c *Config) GetLeaderUUID() types.UUID {
	return c.leaderUuid
}

func (c *Config) GetLeaderUUIDAsString() string {
	return c.leaderUuid.String()
}

func (c *Config) SetLeaderUUID(id types.UUID) {
	c.leaderUuid = id
}

func (c *Config) GetUptimeSeconds() float64 {
	return time.Since(c.uptime).Seconds()
}

func (c *Config) GetMaxLeaderFailures() int {
	return c.maxLeaderFailures
}

func (c *Config) GetPropagationInterval() time.Duration {
	return c.propagationInterval
}

func (c *Config) GetHeartbeatInterval() time.Duration {
	return c.heartbeatInterval
}

func (c *Config) GetHeartbeatTimeout() time.Duration {
	return c.heartbeatTimeout
}

func (c *Config) GetRenewLockInterval() time.Duration {
	return c.renewLockInterval
}

func (c *Config) GetRenewLockTimeout() time.Duration {
	return c.renewLockTimeout
}

func (c *Config) IsLeader() bool {
	return c.leaderUuid == c.id
}

func InitConfig(ctx context.Context) {
	id, err := uuid.Parse(os.Getenv(utils.TowerIdEnv))
	if err != nil {
		log.Fatalf("invalid id %s in env %s: %v", id, utils.TowerIdEnv, err)
	}

	pool, err := pgxpool.New(ctx, os.Getenv(utils.PostgresURIEnv))
	if err != nil {
		log.Fatalf("failed to establish database connection pool: %v", err)
	}

	channel := initRabbitMQ()
	towersQueue := os.Getenv(utils.TowersQueueEnv)
	auditQueue := os.Getenv(utils.AuditQueueEnv)

	dns := os.Getenv(utils.BaseDnsEnv)

	maxLeaderFailures, err := strconv.Atoi(os.Getenv(utils.MaxLeaderFailuresEnv))
	if err != nil {
		log.Fatalf("failed to parse max leader failures env: %v", err)
	}

	pinterval, err := strconv.Atoi(os.Getenv(utils.PropagationIntervalEnv))
	if err != nil {
		log.Fatalf("failed to parse propagation interval env: %v", err)
	}

	propagationInterval := time.Duration(pinterval) * time.Second

	hinterval, err := strconv.Atoi(os.Getenv(utils.HeartbeatIntervalEnv))
	if err != nil {
		log.Fatalf("failed to parse propagation interval env: %v", err)
	}

	heartbeatInterval := time.Duration(hinterval) * time.Second

	htimeout, err := strconv.Atoi(os.Getenv(utils.HeartbeatTimeoutEnv))
	if err != nil {
		log.Fatalf("failed to parse heartbeat timeout env: %v", err)
	}

	heartbeatTimeout := time.Duration(htimeout) * time.Second

	linterval, err := strconv.Atoi(os.Getenv(utils.RenewLockIntervalEnv))
	if err != nil {
		log.Fatalf("failed to parse renew lock interval env: %v", err)
	}

	renewLockInterval := time.Duration(linterval) * time.Second

	ltimeout, err := strconv.Atoi(os.Getenv(utils.RenewLockTimeoutEnv))
	if err != nil {
		log.Fatalf("failed to parse renew lock timeout env: %v", err)
	}

	renewLockTimeout := time.Duration(ltimeout) * time.Second

	Configuration = &Config{
		id:                  types.UUID(id),
		baseDns:             dns,
		towersQueue:         towersQueue,
		auditQueue:          auditQueue,
		db:                  pool,
		rabbitmq:            channel,
		uptime:              time.Now(),
		maxLeaderFailures:   maxLeaderFailures,
		propagationInterval: propagationInterval,
		heartbeatInterval:   heartbeatInterval,
		heartbeatTimeout:    heartbeatTimeout,
		renewLockInterval:   renewLockInterval,
		renewLockTimeout:    renewLockTimeout,
	}
}

func initRabbitMQ() *amqp.Channel {
	var err error
	brokerconn, err = amqp.Dial(os.Getenv(utils.RabbitMQURIEnv))
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
	}

	ch, err := brokerconn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}

	return ch
}

func CloseRabbitMQ() {
	Configuration.GetRabbitMQChannel().Close()
	brokerconn.Close()
}
