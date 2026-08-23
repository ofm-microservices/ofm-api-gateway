package http

import (
	"errors"
	"os"
	"strings"

	service "api-gateway/internal/application"
	gateway "api-gateway/internal/domain"
	"api-gateway/internal/migration"
	"github.com/gofiber/fiber/v2"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
)

// MigrationAuthHandler exposes the registration migration surface under /api/v2.
// It keeps the microservice flow primary and allows legacy fallback only when
// the initial signup cannot be started.
type MigrationAuthHandler interface {
	RegisterRoutes(router fiber.Router)
}

type migrationAuthHandler struct {
	microservice service.RegistrationService
	legacy       migration.LegacyRegistrationClient
	affinity     migration.SessionAffinity
	log          logging.Logger
	legacyWrites bool
}

// NewMigrationAuthHandler constructs the /api/v2 registration migration adapter.
func NewMigrationAuthHandler(microservice service.RegistrationService, legacy migration.LegacyRegistrationClient, affinity migration.SessionAffinity, log logging.Logger) (MigrationAuthHandler, error) {
	if microservice == nil {
		return nil, ErrNilRegistrationService
	}
	if legacy == nil {
		return nil, errors.New("legacy registration client is nil")
	}
	if affinity == nil {
		return nil, errors.New("registration session affinity is nil")
	}
	if log == nil {
		return nil, ErrNilLogger
	}
	return &migrationAuthHandler{microservice: microservice, legacy: legacy, affinity: affinity, legacyWrites: os.Getenv("MIGRATION_REGISTRATION_LEGACY_WRITES_DISABLED") != "true", log: log.With(logging.String("module", "migration-auth-handler"))}, nil
}

func (h *migrationAuthHandler) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/sign-up", h.signUp)
	auth.Post("/sign-up/verify-email", h.verifyEmail)
	auth.Post("/sign-up/complete", h.complete)
}

func (h *migrationAuthHandler) signUp(c *fiber.Ctx) error {
	var req gateway.SignUpRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	result, err := h.microservice.SignUp(c.UserContext(), req)
	if err == nil {
		return c.Status(fiber.StatusAccepted).JSON(result)
	}
	// Recovery commands must retry the owning microservice flow. Falling back
	// to the monolith again here would create a recovery loop and could append
	// another command to the same Kafka partition indefinitely.
	if c.Get("X-Recovery-Replay") == "true" || c.Get("X-Recovery-Loop-Guard") == "true" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "registration service unavailable"})
	}
	if !errors.Is(err, gateway.ErrFailedToStartRegistration) {
		return h.mapRegistrationError(c, err)
	}
	if !h.legacyWrites {
		h.log.Warn("legacy registration write rejected after cutover", logging.Err(err))
		return c.Status(fiber.StatusGone).JSON(fiber.Map{"error": "legacy registration writes are disabled"})
	}

	legacyCtx := migration.WithRecoveryMetadata(c.UserContext(), migration.RecoveryMetadata{
		CommandID: c.Get("X-Command-ID"), CorrelationID: c.Get("X-Correlation-ID"), IdempotencyKey: c.Get("Idempotency-Key"),
	})
	legacyResult, legacyErr := h.legacy.Start(legacyCtx, req)
	if legacyErr != nil {
		return h.mapRegistrationError(c, legacyErr)
	}
	if legacyResult == nil || strings.TrimSpace(legacyResult.SessionID) == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid legacy registration response"})
	}
	// Recovery requests carry their durable command identity in Kafka and must
	// not depend on the optional Redis session-affinity side effect. The
	// monolith outbox is the source of truth for replay; Redis is only needed
	// for interactive legacy verify/complete requests.
	if c.Get("X-Command-ID") == "" && c.Get("X-Recovery-Loop-Guard") != "true" && c.Get("X-Recovery-Replay") != "true" {
		if err := h.affinity.SetLegacy(c.UserContext(), legacyResult.SessionID); err != nil {
			h.log.Error("store legacy registration affinity failed", logging.Err(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
	}
	return c.Status(fiber.StatusAccepted).JSON(legacyResult)
}

func (h *migrationAuthHandler) verifyEmail(c *fiber.Ctx) error {
	var req gateway.VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	if strings.TrimSpace(req.SessionID) == "" || strings.TrimSpace(req.ClientID) == "" || strings.TrimSpace(req.Code) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid verification request"})
	}
	isLegacy, err := h.affinity.IsLegacy(c.UserContext(), req.SessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
	if isLegacy {
		result, err := h.legacy.Verify(c.UserContext(), req)
		if err != nil {
			return h.mapRegistrationError(c, err)
		}
		return c.Status(fiber.StatusAccepted).JSON(result)
	}
	result, err := h.microservice.VerifyEmail(c.UserContext(), req)
	if err != nil {
		return h.mapRegistrationError(c, err)
	}
	return c.Status(fiber.StatusAccepted).JSON(result)
}

func (h *migrationAuthHandler) complete(c *fiber.Ctx) error {
	var req gateway.CompleteRegistrationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	if strings.TrimSpace(req.SessionID) == "" || strings.TrimSpace(req.ClientID) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid completion request"})
	}
	isLegacy, err := h.affinity.IsLegacy(c.UserContext(), req.SessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
	if isLegacy {
		result, err := h.legacy.Complete(c.UserContext(), req)
		if err != nil {
			return h.mapRegistrationError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(result)
	}
	result, err := h.microservice.CompleteRegistration(c.UserContext(), req)
	if err != nil {
		return h.mapRegistrationError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *migrationAuthHandler) mapRegistrationError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, gateway.ErrInvalidEmail):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid email"})
	case errors.Is(err, gateway.ErrInvalidPassword):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid password"})
	case errors.Is(err, gateway.ErrInvalidUsername):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid username"})
	case errors.Is(err, gateway.ErrInvalidSessionID):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid session id"})
	case errors.Is(err, gateway.ErrInvalidClientID):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid client id"})
	case errors.Is(err, gateway.ErrInvalidVerificationCode):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid verification code"})
	case errors.Is(err, gateway.ErrRegistrationAlreadyClaimed):
		return c.Status(fiber.StatusGone).JSON(fiber.Map{"error": "registration tokens already claimed"})
	case errors.Is(err, gateway.ErrRegistrationNotCompleted):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "registration is not completed"})
	}
	var conflict *gateway.RegistrationConflictError
	if errors.As(err, &conflict) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "registration conflict", "state": conflict.State, "username_taken": conflict.UsernameTaken, "email_taken": conflict.EmailTaken})
	}
	h.log.Error("migration registration request failed", logging.Err(err))
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
}
