package appfx

import (
	"api-gateway/config"
	service "api-gateway/internal/application"
	httpserver "api-gateway/internal/presentation/http"
	"github.com/ofm-microservices/ofm-common/pkg/logging"

	"go.uber.org/fx"
)

// HTTPV1Module wires versioned HTTP handlers under the global /v1 prefix.
var HTTPV1Module = fx.Options(
	fx.Provide(ProvideHTTPV1AuthHandler),
	fx.Provide(ProvideHTTPV1GigHandler),
	fx.Provide(ProvideHTTPV1OrderHandler),
	fx.Provide(ProvideHTTPV1ReviewHandler),
	fx.Provide(ProvideHTTPV1SearchHandler),
	fx.Provide(ProvideHTTPV1OnboardingHandler),
	fx.Invoke(InvokeRegisterHTTPV1Routes),
)

// ProvideHTTPV1AuthHandler constructs the versioned auth HTTP handler.
func ProvideHTTPV1AuthHandler(
	registration service.RegistrationService,
	session service.AuthSessionService,
	me service.AuthMeService,
	cfg *config.Config,
	lg logging.Logger,
) (httpserver.AuthHandler, error) {
	return httpserver.NewAuthHandler(registration, session, me, cfg.JWT.AccessSecret, lg)
}

// ProvideHTTPV1GigHandler constructs the versioned gig HTTP handler.
func ProvideHTTPV1GigHandler(
	cfg *config.Config,
	service httpserver.GigService,
	lg logging.Logger,
) (httpserver.GigHandler, error) {
	return httpserver.NewGigHandler(service, cfg.JWT.AccessSecret, lg)
}

// ProvideHTTPV1OrderHandler constructs the versioned order HTTP handler.
func ProvideHTTPV1OrderHandler(
	cfg *config.Config,
	service httpserver.OrderService,
	lg logging.Logger,
) (httpserver.OrderHandler, error) {
	return httpserver.NewOrderHandler(service, cfg.JWT.AccessSecret, lg)
}

// ProvideHTTPV1ReviewHandler constructs the buyer review HTTP handler.
func ProvideHTTPV1ReviewHandler(
	cfg *config.Config,
	service httpserver.ReviewService,
	lg logging.Logger,
) (httpserver.ReviewHandler, error) {
	return httpserver.NewReviewHandler(service, cfg.JWT.AccessSecret, lg)
}

// ProvideHTTPV1SearchHandler constructs the public search HTTP handler.
func ProvideHTTPV1SearchHandler(
	service httpserver.SearchService,
	lg logging.Logger,
) (httpserver.SearchHandler, error) {
	return httpserver.NewSearchHandler(service, lg)
}

// ProvideHTTPV1OnboardingHandler constructs the freelancer onboarding handler.
func ProvideHTTPV1OnboardingHandler(
	cfg *config.Config,
	service httpserver.PaymentOnboardingService,
	lg logging.Logger,
) (httpserver.OnboardingHandler, error) {
	return httpserver.NewOnboardingHandler(service, cfg.JWT.AccessSecret, lg)
}

// InvokeRegisterHTTPV1Routes registers versioned HTTP routes on the server.
func InvokeRegisterHTTPV1Routes(srv httpserver.Server, authHandler httpserver.AuthHandler, gigHandler httpserver.GigHandler, orderHandler httpserver.OrderHandler, reviewHandler httpserver.ReviewHandler, searchHandler httpserver.SearchHandler, onboardingHandler httpserver.OnboardingHandler) {
	srv.App().Get("/v1/search", searchHandler.HandleSearch)

	v1 := srv.App().Group("/v1")
	authHandler.RegisterRoutes(v1)
	gigHandler.RegisterRoutes(v1)
	orderHandler.RegisterRoutes(v1)
	reviewHandler.RegisterRoutes(v1)
	onboardingHandler.RegisterRoutes(v1)
}
