package appfx

import (
	"api-gateway/config"
	"api-gateway/internal/migration"
	httpserver "api-gateway/internal/presentation/http"
	"context"
	"github.com/ofm-microservices/ofm-common/pkg/logging"

	"go.uber.org/fx"
)

// HTTPModule wires the public HTTP server into the FX lifecycle.
var HTTPModule = fx.Options(
	fx.Provide(ProvideHTTPServer),
	fx.Invoke(InvokeRunHTTPServer),
)

// ProvideHTTPServer constructs the HTTP server from gateway config.
func ProvideHTTPServer(cfg *config.Config, lg logging.Logger) (httpserver.Server, error) {
	if cfg.Monolith.BaseURL == "" {
		return httpserver.NewServer(cfg.HTTP, lg)
	}
	if !cfg.Migration.WriteFallback {
		return httpserver.NewServer(cfg.HTTP, lg)
	}
	fallback, err := migration.NewFallbackWriteClient(cfg.Monolith.BaseURL, cfg.Monolith.RateLimitBypassToken, cfg.Monolith.HTTPTimeout)
	if err != nil {
		return nil, err
	}
	return httpserver.NewServer(cfg.HTTP, lg, fallback)
}

// InvokeRunHTTPServer starts and stops the HTTP server with FX lifecycle hooks.
func InvokeRunHTTPServer(lc fx.Lifecycle, srv httpserver.Server, lg logging.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := srv.Start(); err != nil {
					lg.Error("http server stopped with error", logging.Err(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}
