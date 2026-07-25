package config

// MetricsConfig defines the Prometheus scrape listener for api-gateway.
type MetricsConfig struct {
	Enabled bool   `env:"METRICS_ENABLED" envDefault:"true"`
	Host    string `env:"METRICS_HOST" envDefault:"0.0.0.0"`
	Port    int    `env:"METRICS_PORT" envDefault:"9606"`
	Path    string `env:"METRICS_PATH" envDefault:"/metrics"`
}
