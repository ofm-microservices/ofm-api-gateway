package http

import (
	gateway "api-gateway/internal/domain"
	"errors"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"

	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	service RegistrationService
	log     logging.Logger
}

// NewAuthHandler constructs the auth HTTP handler group for signup requests.
func NewAuthHandler(service RegistrationService, log logging.Logger) (AuthHandler, error) {
	if service == nil {
		return nil, ErrNilRegistrationService
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	return &authHandler{
		service: service,
		log:     log.With(logging.String("module", "http-auth-handler")),
	}, nil
}

// RegisterRoutes mounts auth routes under the router it receives.
func (h *authHandler) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Post("/sign-up", h.HandleSignUp)
}

// HandleSignUp parses the public signup payload and starts registration.
func (h *authHandler) HandleSignUp(c *fiber.Ctx) error {
	var req gateway.SignUpRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	result, err := h.service.SignUp(c.UserContext(), req)
	if err != nil {
		return h.MapSignUpError(c, err)
	}

	h.log.Info("sign up request accepted")

	return c.Status(fiber.StatusAccepted).JSON(result)
}

// MapSignUpError translates signup failures into stable HTTP responses.
func (h *authHandler) MapSignUpError(c *fiber.Ctx, err error) error {
	var conflictErr *gateway.RegistrationConflictError

	switch err {
	case gateway.ErrInvalidEmail:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid email"})
	case gateway.ErrInvalidPassword:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid password"})
	case gateway.ErrInvalidUsername:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid username"})
	default:
		if errors.As(err, &conflictErr) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":          "registration conflict",
				"state":          conflictErr.State,
				"username_taken": conflictErr.UsernameTaken,
				"email_taken":    conflictErr.EmailTaken,
			})
		}
		h.log.Error("request failed", logging.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}
