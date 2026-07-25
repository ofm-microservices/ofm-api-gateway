package config

// ReviewServiceConfig configures the outbound gRPC client used to submit
// buyer reviews.
type ReviewServiceConfig struct {
	Address string `env:"REVIEW_SERVICE_ADDRESS" envDefault:"127.0.0.1:9510"`
}
