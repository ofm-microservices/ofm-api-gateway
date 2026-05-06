package http

import (
	service "api-gateway/internal/application"
	"context"

	"github.com/gofiber/fiber/v2"
)

// RegistrationService aliases the application regis
// the HTTP layer.
type RegistrationService = service.RegistrationService

// AuthHandler exposes the auth HTTP routes owned by api-gateway.
type AuthHandler interface {
	RegisterRoutes(router fiber.Router)
	HandleSignUp(c *fiber.Ctx) error
	HandleVerifyEmail(c *fiber.Ctx) error
	HandleCompleteRegistration(c *fiber.Ctx) error
}

// Server exposes the api-gateway HTTP server lifecycle.
type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
	App() *fiber.App
}
