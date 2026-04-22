package appfx

import (
	"api-gateway/config"

	"go.uber.org/fx"
)

// ConfigModule provides configuration loading for api-gateway.
var ConfigModule = fx.Options(
	fx.Provide(ProvideConfig),
)

// ProvideConfig loads and validates the api-gateway configuration.
func ProvideConfig() (*config.Config, error) {
	return config.Load()
}
