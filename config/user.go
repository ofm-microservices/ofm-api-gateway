package config

// UserServiceConfig configures the outbound gRPC client used for public user
// profile lookups.
type UserServiceConfig struct {
	Address string `env:"USER_SERVICE_ADDRESS" envDefault:"127.0.0.1:9502"`
}
