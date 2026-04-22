package appfx

import (
	service "api-gateway/internal/application"
	httpserver "api-gateway/internal/presentation/http"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"

	"go.uber.org/fx"
)

// HTTPV1Module wires versioned HTTP handlers under the global /v1 prefix.
var HTTPV1Module = fx.Options(
	fx.Provide(ProvideHTTPV1AuthHandler),
	fx.Invoke(InvokeRegisterHTTPV1Routes),
)

// ProvideHTTPV1AuthHandler constructs the versioned auth HTTP handler.
func ProvideHTTPV1AuthHandler(
	service service.RegistrationService,
	lg logging.Logger,
) (httpserver.AuthHandler, error) {
	return httpserver.NewAuthHandler(service, lg)
}

// InvokeRegisterHTTPV1Routes registers versioned HTTP routes on the server.
func InvokeRegisterHTTPV1Routes(srv httpserver.Server, authHandler httpserver.AuthHandler) {
	v1 := srv.App().Group("/v1")
	authHandler.RegisterRoutes(v1)
}
