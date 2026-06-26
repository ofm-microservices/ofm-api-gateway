package config

// ChatServiceConfig configures the outbound gRPC client used for order chat flows.
type ChatServiceConfig struct {
	Address string `env:"CHAT_SERVICE_ADDRESS"`
}
