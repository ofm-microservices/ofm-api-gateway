package config

// NATSConfig defines the broker connection settings used by api-gateway.
type NATSConfig struct {
	URL                   string `env:"NATS_URL" envDefault:"nats://127.0.0.1:4222"`
	User                  string `env:"NATS_USER"`
	Password              string `env:"NATS_PASSWORD"`
	SagaCreateAuthSubject string `env:"NATS_SUBJECT_SAGA_CREATE_AUTH" envDefault:"saga.auth.create"`
	OrderSagaStartSubject string `env:"NATS_SUBJECT_ORDER_SAGA_START" envDefault:"order.saga.start"`
}
