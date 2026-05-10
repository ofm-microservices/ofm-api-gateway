package config

// JWTConfig configures JWT verification in the api gateway.
type JWTConfig struct {
	Secret string `env:"JWT_SECRET" envDefault:"local-dev-secret-change-me"`
}
