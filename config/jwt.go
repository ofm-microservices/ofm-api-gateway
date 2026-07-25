package config

// JWTConfig configures JWT verification in the api gateway.
type JWTConfig struct {
	AccessSecret string `env:"JWT_ACCESS_SECRET" envDefault:"local-dev-access-secret-change-me"`
}
