package appfx

import (
	"api-gateway/config"
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
	return httpserver.NewServer(cfg.HTTP, lg)
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
