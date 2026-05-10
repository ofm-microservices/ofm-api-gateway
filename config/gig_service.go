package config

// GigServiceConfig configures the internal gRPC client used to manage gig
// drafts.
type GigServiceConfig struct {
	Address string `env:"GIG_SERVICE_ADDRESS" envDefault:"127.0.0.1:9093"`
}
