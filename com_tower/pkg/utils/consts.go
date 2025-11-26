package utils

const (
	// envs
	TowerIdEnv              = "TOWER_ID"
	PortEnv                 = "PORT"
	PostgresURIEnv          = "POSTGRES_URI"
	RabbitMQURIEnv          = "RABBITMQ_URI"
	TowersQueueEnv          = "TOWERS_QUEUE"
	AuditQueueEnv           = "AUDIT_QUEUE"
	MaxLeaderFailuresEnv    = "MAX_LEADER_FAILURES"
	MaxStructureFailuresEnv = "MAX_STRUCTURE_FAILURES"
	PropagationIntervalEnv  = "PROPAGATION_INTERVAL"
	HeartbeatIntervalEnv    = "HEARTBEAT_INTERVAL"
	HeartbeatTimeoutEnv     = "HEARTBEAT_TIMEOUT"
	RenewLockIntervalEnv    = "RENEW_LOCK_INTERVAL"
	RenewLockTimeoutEnv     = "RENEW_LOCK_TIMEOUT"
	BaseDnsEnv              = "BASE_DNS"
	EmailHostEnv            = "EMAIL_HOST"
	EmailPortEnv            = "EMAIL_PORT"
	EmailUserEnv            = "EMAIL_USER"
	EmailPasswordEnv        = "EMAIL_PASSWORD"
	EmailRecipientsEnv      = "EMAIL_RECIPIENTS"

	// email templates
	EmailSubjectTemplate = "[CRITICAL] %s %s down!"
	EmailBodyTemplate = "Alert!\nTower %s has identified that %s %s is down!\nPlease check the status of the structure right now!"
)
