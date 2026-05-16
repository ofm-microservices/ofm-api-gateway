package config

// PaymentServiceConfig configures the internal gRPC client used for payment
// onboarding flows.
type PaymentServiceConfig struct {
	Address string `env:"PAYMENT_SERVICE_ADDRESS,required"`
}
