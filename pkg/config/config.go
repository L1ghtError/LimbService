package config

// values that represents all supported params by application
const (
	ClientKey      = "Client"
	SecretKey      = "Secret"
	RedirectURLKey = "Redirect_URL"
	AppPortKey     = "APP_PORT"
	AppHostKey     = "APP_HOST"

	DbHostKey     = "DB_HOST"
	DbPortKey     = "DB_PORT"
	DbPasswordKey = "DB_PASSWORD"
	DbNameKey     = "DB_NAME"

	JwtAccessSecretKey  = "JWT_ACCESS_SECRET"
	JwtRefreshSecretKey = "JWT_REFRESH_SECRET"

	AmqpUserKey   = "AMQP_USER"
	AmqpPasswdKey = "AMQP_PASSWD"
	AmqpHostKey   = "AMQP_HOST"
	AmqpPortKey   = "AMQP_PORT"
)

type Config interface {
	// Init loads config
	Init(filenames ...string) error

	// Get receive spesific value
	Get(key string) string
}
