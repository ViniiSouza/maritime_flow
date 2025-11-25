package utils

const (
	// envs
	TowerIdEnv             = "TOWER_ID"
	PortEnv                = "PORT"
	PostgresURIEnv         = "POSTGRES_URI"
	RabbitMQURIEnv         = "RABBITMQ_URI"
	TowersQueueEnv         = "TOWERS_QUEUE"
	AuditQueueEnv          = "AUDIT_QUEUE"
	MaxLeaderFailuresEnv   = "MAX_LEADER_FAILURES"
	PropagationIntervalEnv = "PROPAGATION_INTERVAL"
	HeartbeatIntervalEnv   = "HEARTBEAT_INTERVAL"
	HeartbeatTimeoutEnv    = "HEARTBEAT_TIMEOUT"
	RenewLockIntervalEnv   = "RENEW_LOCK_INTERVAL"
	RenewLockTimeoutEnv    = "RENEW_LOCK_TIMEOUT"
	BaseDnsEnv             = "BASE_DNS"
)
