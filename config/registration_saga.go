package config

// RegistrationSagaConfig configures the internal gRPC client used to start the
// registration saga.
type RegistrationSagaConfig struct {
	Address string `env:"REGISTRATION_SAGA_ADDRESS" envDefault:"127.0.0.1:9500"`
}
