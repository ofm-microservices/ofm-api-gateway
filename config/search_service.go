package config

// SearchServiceConfig configures the outbound search-service gRPC client.
type SearchServiceConfig struct {
	Address string `env:"SEARCH_SERVICE_ADDRESS" envDefault:"127.0.0.1:9511"`
}
