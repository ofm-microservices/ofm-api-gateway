package http

import (
	service "api-gateway/internal/application"
	gateway "api-gateway/internal/domain"
	"context"

	"github.com/gofiber/fiber/v2"
)

// RegistrationService aliases the application regis
// the HTTP layer.
type RegistrationService = service.RegistrationService

// AuthSessionService aliases the sign-in application boundary used by the HTTP layer.
type AuthSessionService = service.AuthSessionService

// AuthHandler exposes the auth HTTP routes owned by api-gateway.
type AuthHandler interface {
	RegisterRoutes(router fiber.Router)
	HandleSignUp(c *fiber.Ctx) error
	HandleSignIn(c *fiber.Ctx) error
	HandleVerifyEmail(c *fiber.Ctx) error
	HandleCompleteRegistration(c *fiber.Ctx) error
}

// GigService aliases the application gig orchestration contract used by the
// HTTP layer.
type GigService = service.GigService

// OrderService aliases the application order orchestration contract used by
// the HTTP layer.
type OrderService = service.OrderService

// ReviewService aliases the application review orchestration contract used by
// the HTTP layer.
type ReviewService = service.ReviewService

// SearchService aliases the application search orchestration contract used by
// the HTTP layer.
type SearchService = service.SearchService

// GigPrincipalResolver validates bearer tokens and extracts the gig owner id.
type GigPrincipalResolver interface {
	Middleware() fiber.Handler
	FreelancerID(c *fiber.Ctx) (string, error)
	Email(c *fiber.Ctx) (string, error)
}

// GigHandler exposes the gig draft HTTP routes owned by api-gateway.
type GigHandler interface {
	RegisterRoutes(router fiber.Router)
	HandleCreateDraft(c *fiber.Ctx) error
	HandleUpdateBasicInfo(c *fiber.Ctx) error
	HandleReplacePackages(c *fiber.Ctx) error
	HandleReplaceQuestions(c *fiber.Ctx) error
	HandleReplaceMedia(c *fiber.Ctx) error
	HandleGetDraft(c *fiber.Ctx) error
	HandleGetBySlug(c *fiber.Ctx) error
	HandlePublish(c *fiber.Ctx) error
}

// OrderHandler exposes the create-order HTTP route owned by api-gateway.
type OrderHandler interface {
	RegisterRoutes(router fiber.Router)
	HandleCreateOrder(c *fiber.Ctx) error
	HandleConfirmOrder(c *fiber.Ctx) error
	HandleDeliverOrder(c *fiber.Ctx) error
	HandleAcceptDelivery(c *fiber.Ctx) error
	HandleRequestRevision(c *fiber.Ctx) error
	HandleOpenDispute(c *fiber.Ctx) error
}

// ReviewHandler exposes the buyer review HTTP routes owned by api-gateway.
type ReviewHandler interface {
	RegisterRoutes(router fiber.Router)
	HandleCreateReview(c *fiber.Ctx) error
}

// SearchHandler exposes the public search HTTP route owned by api-gateway.
type SearchHandler interface {
	RegisterRoutes(router fiber.Router)
	HandleSearch(c *fiber.Ctx) error
}

// PaymentOnboardingService aliases the onboarding application boundary.
type PaymentOnboardingService = service.PaymentOnboardingService

// OnboardingHandler exposes the freelancer onboarding HTTP route group.
type OnboardingHandler interface {
	RegisterRoutes(router fiber.Router)
	HandleStartFreelancerOnboarding(c *fiber.Ctx) error
}

// GigDraftRequest aliases the public draft request payload.
type GigDraftRequest = gateway.CreateGigDraftRequest

// GigBasicInfoRequest aliases the public basic-info update payload.
type GigBasicInfoRequest = gateway.UpdateGigBasicInfoRequest

// GigPackagesRequest aliases the public package replacement payload.
type GigPackagesRequest = gateway.ReplaceGigPackagesRequest

// GigQuestionsRequest aliases the public question replacement payload.
type GigQuestionsRequest = gateway.ReplaceGigQuestionsRequest

// GigMediaRequest aliases the public media replacement payload.
type GigMediaRequest = gateway.ReplaceGigMediaRequest

// GigDraftLookupRequest aliases the public draft lookup payload.
type GigDraftLookupRequest = gateway.GetGigDraftRequest

// GigPublishRequest aliases the public publish payload.
type GigPublishRequest = gateway.PublishGigRequest

// Server exposes the api-gateway HTTP server lifecycle.
type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
	App() *fiber.App
}
