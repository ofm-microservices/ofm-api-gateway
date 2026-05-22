package config

// OrderSagaConfig configures the outbound gRPC client used to orchestrate the
// public order checkout flow.
type OrderSagaConfig struct {
	Address string `env:"ORDER_SAGA_ADDRESS" envDefault:"127.0.0.1:9507"`
}
