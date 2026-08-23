package appfx

import (
	httpserver "api-gateway/internal/presentation/http"
	"go.uber.org/fx"
)

// HTTPV2Module wires the v2 API under the canonical /api/v2 namespace.
var HTTPV2Module = fx.Options(
	fx.Invoke(InvokeRegisterHTTPV2Routes),
)

// InvokeRegisterHTTPV2Routes registers the only public v2 surface under
// /api/v2.
func InvokeRegisterHTTPV2Routes(
	srv httpserver.Server,
	auth httpserver.MigrationAuthHandler,
	authHandler httpserver.AuthHandler,
	user httpserver.MigrationUserHandler,
	userOrder httpserver.UserOrderHandler,
	chat httpserver.ChatHandler,
	gig httpserver.GigHandler,
	order httpserver.OrderHandler,
	review httpserver.ReviewHandler,
	search httpserver.SearchHandler,
	onboarding httpserver.OnboardingHandler,
) {
	v2 := srv.App().Group("/api/v2")
	auth.RegisterRoutes(v2)
	authRoutes := v2.Group("/auth")
	authRoutes.Post("/sign-in", authHandler.HandleSignIn)
	authRoutes.Post("/refresh", authHandler.HandleRefresh)
	authRoutes.Post("/sign-out", authHandler.HandleSignOut)
	authRoutes.Get("/me", authHandler.HandleMe)
	user.RegisterMigrationRoutes(v2)
	userOrder.RegisterRoutes(v2)
	chat.RegisterRoutes(v2)
	gig.RegisterRoutes(v2)
	order.RegisterRoutes(v2)
	review.RegisterRoutes(v2)
	search.RegisterRoutes(v2)
	onboarding.RegisterRoutes(v2)
}
