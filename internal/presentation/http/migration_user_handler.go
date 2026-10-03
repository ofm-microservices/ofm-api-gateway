package http

import (
	"errors"
	"net/url"
	"strings"

	service "api-gateway/internal/application"
	gateway "api-gateway/internal/domain"
	"api-gateway/internal/migration"
	"github.com/gofiber/fiber/v2"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
)

// MigrationUserHandler exposes the v2 user read route with controlled legacy fallback.
type MigrationUserHandler interface{ RegisterMigrationRoutes(fiber.Router) }

type migrationUserHandler struct {
	service service.UserProfileService
	legacy  migration.LegacyReadClient
	log     logging.Logger
}

// NewMigrationUserHandler constructs the user migration adapter.
func NewMigrationUserHandler(svc service.UserProfileService, legacy migration.LegacyReadClient, log logging.Logger) (MigrationUserHandler, error) {
	if svc == nil {
		return nil, ErrNilUserProfileService
	}
	if legacy == nil {
		return nil, errors.New("legacy read client is nil")
	}
	if log == nil {
		return nil, ErrNilLogger
	}
	return &migrationUserHandler{service: svc, legacy: legacy, log: log.With(logging.String("module", "migration-user-handler"))}, nil
}

func (h *migrationUserHandler) RegisterMigrationRoutes(router fiber.Router) {
	router.Get("/users/:username", h.get)
}

func (h *migrationUserHandler) get(c *fiber.Ctx) error {
	username, err := url.PathUnescape(c.Params("username"))
	if err != nil || strings.TrimSpace(username) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": gateway.ErrInvalidUsername.Error()})
	}
	req := gateway.GetUserProfileRequest{Username: strings.TrimSpace(username), GigsCursor: c.Query("gigs_cursor"), ReviewsCursor: c.Query("reviews_cursor")}
	result, err := h.service.GetUserProfile(c.UserContext(), req)
	if err == nil {
		return c.Status(fiber.StatusOK).JSON(result)
	}
	if errors.Is(err, gateway.ErrInvalidUsername) || errors.Is(err, gateway.ErrUserNotFound) {
		return mapMigrationUserError(c, err)
	}
	raw, status, legacyErr := h.legacy.UserProfile(c.UserContext(), req.Username, req.GigsCursor, req.ReviewsCursor)
	if legacyErr != nil {
		h.log.Error("legacy user read fallback failed", logging.Err(legacyErr))
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "user service unavailable"})
	}
	if status >= 200 && status < 300 {
		c.Status(status)
		return c.Send(raw)
	}
	if status == fiber.StatusNotFound {
		return c.Status(status).Send(raw)
	}
	return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "user service unavailable"})
}

func mapMigrationUserError(c *fiber.Ctx, err error) error {
	if errors.Is(err, gateway.ErrInvalidUsername) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
}
