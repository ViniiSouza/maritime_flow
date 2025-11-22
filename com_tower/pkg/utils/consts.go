package utils

const (
	// envs
	TowerIdEnv             = "TOWER_ID"
	PortEnv                = "PORT"
	PostgresURIEnv         = "POSTGRES_URI"
	RabbitMQURIEnv         = "RABBITMQ_URI"
	TowersQueueEnv         = "TOWERS_QUEUE"
	AuditQueueEnv          = "AUDIT_QUEUE"
	PropagationIntervalEnv = "PROPAGATION_INTERVAL"
	HeartbeatIntervalEnv   = "HEARTBEAT_INTERVAL"
	HeartbeatTimeoutEnv    = "HEARTBEAT_TIMEOUT"
	BaseDnsEnv             = "BASE_DNS"

	// paths
	TowersPropagationPath     = "/towers"
	StructuresPropagationPath = "/structures"
)
