package config

// HTTPConfig defines the public HTTP listener for api-gateway.
type HTTPConfig struct {
	Host           string `env:"HTTP_HOST" envDefault:"0.0.0.0"`
	Port           int    `env:"HTTP_PORT" envDefault:"8080"`
	Environment    string `env:"APP_ENV" envDefault:"local"`
	FaultInjection FaultInjectionConfig
}
