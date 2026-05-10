package appfx

import (
	"api-gateway/config"
	service "api-gateway/internal/application"
	httpserver "api-gateway/internal/presentation/http"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"

	"go.uber.org/fx"
)

// HTTPV1Module wires versioned HTTP handlers under the global /v1 prefix.
var HTTPV1Module = fx.Options(
	fx.Provide(ProvideHTTPV1AuthHandler),
	fx.Provide(ProvideHTTPV1GigHandler),
	fx.Invoke(InvokeRegisterHTTPV1Routes),
)

// ProvideHTTPV1AuthHandler constructs the versioned auth HTTP handler.
func ProvideHTTPV1AuthHandler(
	service service.RegistrationService,
	lg logging.Logger,
) (httpserver.AuthHandler, error) {
	return httpserver.NewAuthHandler(service, lg)
}

// ProvideHTTPV1GigHandler constructs the versioned gig HTTP handler.
func ProvideHTTPV1GigHandler(
	cfg *config.Config,
	service httpserver.GigService,
	lg logging.Logger,
) (httpserver.GigHandler, error) {
	return httpserver.NewGigHandler(service, cfg.JWT.Secret, lg)
}

// InvokeRegisterHTTPV1Routes registers versioned HTTP routes on the server.
func InvokeRegisterHTTPV1Routes(srv httpserver.Server, authHandler httpserver.AuthHandler, gigHandler httpserver.GigHandler) {
	v1 := srv.App().Group("/v1")
	authHandler.RegisterRoutes(v1)
	gigHandler.RegisterRoutes(v1)
}
