package config

// AuthServiceConfig configures the internal gRPC client used to issue final
// registration login tokens.
type AuthServiceConfig struct {
	Address string `env:"AUTH_SERVICE_ADDRESS" envDefault:"127.0.0.1:9501"`
}
