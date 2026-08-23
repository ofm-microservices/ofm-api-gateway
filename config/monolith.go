package config

import "time"

// MonolithConfig defines the legacy HTTP endpoint and session-affinity store
// used by the registration migration router.
type MonolithConfig struct {
	BaseURL              string        `env:"MONOLITH_BASE_URL" envDefault:"http://127.0.0.1:8000"`
	RateLimitBypassToken string        `env:"MONOLITH_RATE_LIMIT_BYPASS_TOKEN"`
	HTTPTimeout          time.Duration `env:"MONOLITH_HTTP_TIMEOUT" envDefault:"5s"`
	AffinityTTL          time.Duration `env:"MIGRATION_SESSION_AFFINITY_TTL" envDefault:"20m"`
	RedisHost            string        `env:"MIGRATION_REDIS_HOST" envDefault:"127.0.0.1"`
	RedisPort            int           `env:"MIGRATION_REDIS_PORT" envDefault:"6379"`
	RedisPassword        string        `env:"MIGRATION_REDIS_PASSWORD"`
	RedisDB              int           `env:"MIGRATION_REDIS_DB" envDefault:"0"`
}
