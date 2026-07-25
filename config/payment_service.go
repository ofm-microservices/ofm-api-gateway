package config

// PaymentServiceConfig configures the internal gRPC client used for payment
// onboarding flows.
type PaymentServiceConfig struct {
	Address string `env:"PAYMENT_SERVICE_ADDRESS" envDefault:"127.0.0.1:9506"`
}
