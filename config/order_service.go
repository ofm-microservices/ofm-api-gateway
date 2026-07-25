package config

// OrderServiceConfig configures the outbound gRPC client used to resolve user-scoped order previews.
type OrderServiceConfig struct {
	Address string `env:"ORDER_SERVICE_ADDRESS"`
}
