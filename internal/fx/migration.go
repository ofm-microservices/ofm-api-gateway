package appfx

import (
	"api-gateway/config"
	service "api-gateway/internal/application"
	"api-gateway/internal/migration"
	httpserver "api-gateway/internal/presentation/http"
	"context"
	"fmt"

	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// MigrationModule wires the legacy registration adapter and session affinity.
var MigrationModule = fx.Options(
	fx.Invoke(InvokeRunRecoveryConsumer),
	fx.Provide(ProvideMigrationRedis),
	fx.Provide(ProvideSessionAffinity),
	fx.Provide(ProvideLegacyRegistrationClient),
	fx.Provide(ProvideLegacyReadClient),
	fx.Provide(ProvideLegacySearchClient),
	fx.Provide(ProvideRecoveryConsumer),
	fx.Provide(ProvideMigrationAuthHandler),
	fx.Provide(ProvideMigrationUserHandler),
)

// ProvideRecoveryConsumer constructs the Kafka consumer that replays fallback
// commands through the gateway's normal HTTP and gRPC path.
func ProvideRecoveryConsumer(cfg *config.Config) (migration.RecoveryConsumer, error) {
	return migration.NewRecoveryConsumer(cfg.Migration.RecoveryBrokers, cfg.Migration.RecoveryTopic, cfg.Migration.RecoveryGroup, cfg.Migration.RecoveryDLQ, cfg.Migration.RecoveryBaseURL, cfg.Monolith.HTTPTimeout, cfg.Migration.RecoveryCompleted)
}

// InvokeRunRecoveryConsumer starts and stops fallback command replay with the
// gateway lifecycle.
func InvokeRunRecoveryConsumer(lc fx.Lifecycle, consumer migration.RecoveryConsumer, lg logging.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := consumer.Run(context.Background()); err != nil {
					lg.Error("recovery consumer stopped", logging.Err(err))
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error { return consumer.Close() },
	})
}

// ProvideMigrationRedis constructs the Redis client used for session affinity.
func ProvideMigrationRedis(lc fx.Lifecycle, cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", cfg.Monolith.RedisHost, cfg.Monolith.RedisPort), Password: cfg.Monolith.RedisPassword, DB: cfg.Monolith.RedisDB})
	lc.Append(fx.Hook{OnStop: func(context.Context) error { return client.Close() }})
	return client, nil
}

// ProvideSessionAffinity constructs the Redis-backed legacy session owner store.
func ProvideSessionAffinity(client *redis.Client, cfg *config.Config) (migration.SessionAffinity, error) {
	return migration.NewRedisSessionAffinity(client, cfg.Monolith.AffinityTTL)
}

// ProvideLegacyRegistrationClient constructs the monolith registration adapter.
func ProvideLegacyRegistrationClient(cfg *config.Config) (migration.LegacyRegistrationClient, error) {
	return migration.NewLegacyRegistrationClient(cfg.Monolith.BaseURL, cfg.Monolith.HTTPTimeout)
}

// ProvideLegacyReadClient constructs the read-only fallback adapter.
func ProvideLegacyReadClient(cfg *config.Config) (migration.LegacyReadClient, error) {
	return migration.NewLegacyReadClient(cfg.Monolith.BaseURL, cfg.Monolith.HTTPTimeout)
}

// ProvideLegacySearchClient constructs the monolith-backed search adapter.
func ProvideLegacySearchClient(cfg *config.Config) (migration.LegacySearchClient, error) {
	return migration.NewLegacySearchClient(cfg.Monolith.BaseURL, cfg.Monolith.HTTPTimeout)
}

// ProvideMigrationAuthHandler constructs the /api/v2 registration handler.
func ProvideMigrationAuthHandler(reg service.RegistrationService, legacy migration.LegacyRegistrationClient, affinity migration.SessionAffinity, lg logging.Logger) (httpserver.MigrationAuthHandler, error) {
	return httpserver.NewMigrationAuthHandler(reg, legacy, affinity, lg)
}

// ProvideMigrationUserHandler constructs the v2 user read migration handler.
func ProvideMigrationUserHandler(svc service.UserProfileService, legacy migration.LegacyReadClient, lg logging.Logger) (httpserver.MigrationUserHandler, error) {
	return httpserver.NewMigrationUserHandler(svc, legacy, lg)
}
