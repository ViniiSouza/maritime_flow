package config

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/utils"
)

var brokerconn *amqp.Connection
var Configuration *Config

type Config struct {
	id          types.UUID
	baseDns     string
	leaderUuid  types.UUID
	towersQueue string
	auditQueue  string

	db       *pgx.Conn
	rabbitmq *amqp.Channel

	propagationInterval time.Duration
	heartbeatInterval   time.Duration
	heartbeatTimeout    time.Duration
}

func (c *Config) GetId() types.UUID {
	return c.id
}

func (c *Config) GetIdAsString() string {
	return c.id.String()
}

func (c *Config) GetDBConn() *pgx.Conn {
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

func (c *Config) SetLeaderUUID(id types.UUID) {
	c.leaderUuid = id
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

func InitConfig(ctx context.Context) {
	id, err := uuid.Parse(os.Getenv(utils.TowerIdEnv))
	if err != nil {
		log.Fatalf("invalid id %s in env %s: %v", id, utils.TowerIdEnv, err)
	}

	conn, err := pgx.Connect(ctx, os.Getenv(utils.PostgresURIEnv))
	if err != nil {
		log.Fatalf("failed to establish database connection: %v", err)
	}

	channel := initRabbitMQ()
	towersQueue := os.Getenv(utils.TowersQueueEnv)
	auditQueue := os.Getenv(utils.AuditQueueEnv)

	dns := os.Getenv(utils.BaseDnsEnv)

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

	timeout, err := strconv.Atoi(os.Getenv(utils.HeartbeatTimeoutEnv))
	if err != nil {
		log.Fatalf("failed to parse heartbeat timeout env: %v", err)
	}

	heartbeatTimeout := time.Duration(timeout) * time.Second
	Configuration = &Config{
		id:                  types.UUID(id),
		baseDns:             dns,
		towersQueue:         towersQueue,
		auditQueue:          auditQueue,
		db:                  conn,
		rabbitmq:            channel,
		propagationInterval: propagationInterval,
		heartbeatInterval:   heartbeatInterval,
		heartbeatTimeout:    heartbeatTimeout,
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
