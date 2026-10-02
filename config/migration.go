package config

// MigrationConfig controls strangler-fig routing without coupling handlers to
// legacy storage or service implementations.
type MigrationConfig struct {
	AuthMode         string `env:"MIGRATION_AUTH_MODE" envDefault:"microservice"`
	UserMode         string `env:"MIGRATION_USER_MODE" envDefault:"microservice"`
	RegistrationMode string `env:"MIGRATION_REGISTRATION_MODE" envDefault:"microservice"`
	OnboardingMode   string `env:"MIGRATION_ONBOARDING_MODE" envDefault:"monolith"`
	GigMode          string `env:"MIGRATION_GIG_MODE" envDefault:"monolith"`
	OrderMode        string `env:"MIGRATION_ORDER_MODE" envDefault:"monolith"`
	PaymentMode      string `env:"MIGRATION_PAYMENT_MODE" envDefault:"monolith"`
	ReviewMode       string `env:"MIGRATION_REVIEW_MODE" envDefault:"monolith"`
	SearchMode       string `env:"MIGRATION_SEARCH_MODE" envDefault:"monolith"`
	AuthReadFallback bool   `env:"MIGRATION_AUTH_READ_FALLBACK" envDefault:"false"`
	UserReadFallback bool   `env:"MIGRATION_USER_READ_FALLBACK" envDefault:"true"`
	WriteFallback    bool   `env:"MIGRATION_WRITE_FALLBACK" envDefault:"true"`
}
