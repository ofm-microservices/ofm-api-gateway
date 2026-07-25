package http

import (
	gateway "api-gateway/internal/domain"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"net/url"
	"strings"
	"time"
)

type userHandler struct {
	service UserProfileService
	log     logging.Logger
}

// NewUserHandler constructs the public user-profile HTTP handler group.
func NewUserHandler(service UserProfileService, log logging.Logger) (UserHandler, error) {
	if service == nil {
		return nil, ErrNilUserProfileService
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	return &userHandler{
		service: service,
		log:     log.With(logging.String("module", "http-user-handler")),
	}, nil
}

func (h *userHandler) RegisterRoutes(router fiber.Router) {
	users := router.Group("/users/:username")
	users.Get("", h.HandleGetByUsername)
}

func (h *userHandler) HandleGetByUsername(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)
	username, err := url.PathUnescape(c.Params("username"))
	if err != nil {
		return h.mapUserError(c, gateway.ErrInvalidUsername)
	}
	username = strings.TrimSpace(username)
	if username == "" {
		return h.mapUserError(c, gateway.ErrInvalidUsername)
	}
	req := gateway.GetUserProfileRequest{
		Username:      username,
		GigsCursor:    c.Query("gigs_cursor"),
		ReviewsCursor: c.Query("reviews_cursor"),
	}

	result, err := h.service.GetUserProfile(c.UserContext(), req)
	if err != nil {
		return h.mapUserError(c, err)
	}

	log.Info("user profile request accepted",
		logging.Operation("http.user.profile"),
		logging.DurationMS(time.Since(started)),
		logging.String("username", req.Username),
	)
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *userHandler) mapUserError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, gateway.ErrInvalidUsername):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}
