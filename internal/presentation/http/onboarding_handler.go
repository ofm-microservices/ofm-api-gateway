package http

import (
	gateway "api-gateway/internal/domain"
	"errors"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
)

type onboardingHandler struct {
	service PaymentOnboardingService
	auth    GigPrincipalResolver
	log     logging.Logger
}

// NewOnboardingHandler constructs the freelancer onboarding HTTP handler group.
func NewOnboardingHandler(service PaymentOnboardingService, jwtSecret string, log logging.Logger) (OnboardingHandler, error) {
	if service == nil {
		return nil, ErrNilPaymentOnboardingService
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	auth, err := newJWTPrincipalResolver(jwtSecret, log)
	if err != nil {
		return nil, err
	}

	return &onboardingHandler{
		service: service,
		auth:    auth,
		log:     log.With(logging.String("module", "http-onboarding-handler")),
	}, nil
}

// RegisterRoutes mounts onboarding routes under the provided router.
func (h *onboardingHandler) RegisterRoutes(router fiber.Router) {
	freelancer := router.Group("/freelancer")
	onboarding := freelancer.Group("/onboarding")
	onboarding.Use(h.auth.Middleware())
	onboarding.Post("/start", h.HandleStartFreelancerOnboarding)
}

// HandleStartFreelancerOnboarding starts Stripe Connect onboarding.
func (h *onboardingHandler) HandleStartFreelancerOnboarding(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)
	userID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	result, err := h.service.StartFreelancerOnboarding(c.UserContext(), gateway.StartFreelancerOnboardingRequest{
		UserID:  userID,
		Country: "",
	})
	if err != nil {
		return h.mapOnboardingError(c, err)
	}

	log.Info("freelancer onboarding request accepted",
		logging.Operation("http.onboarding.start"),
		logging.DurationMS(time.Since(started)),
		logging.String("user_id", userID),
		logging.String("stripe_account_id", result.StripeAccountID),
	)
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *onboardingHandler) mapOnboardingError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, gateway.ErrInvalidFreelancerOnboarding):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	default:
		h.log.Error("request failed", logging.Err(err))
		if os.Getenv("APP_ENV") == "local" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}
