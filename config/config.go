package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config groups the full api-gateway runtime configuration.
type Config struct {
	App              AppConfig
	HTTP             HTTPConfig
	RegistrationSaga RegistrationSagaConfig
	AuthService      AuthServiceConfig
}

// Load reads environment variables into Config and applies defaults.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, WrapParseEnvConfigError(err)
	}

	return cfg, nil
}
