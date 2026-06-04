package http

import (
	gateway "api-gateway/internal/domain"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
)

type userOrderHandler struct {
	service OrderPreviewService
	auth    GigPrincipalResolver
	log     logging.Logger
}

// NewUserOrderHandler constructs the authenticated user-scoped order preview HTTP handler group.
func NewUserOrderHandler(service OrderPreviewService, jwtSecret string, log logging.Logger) (UserOrderHandler, error) {
	if service == nil {
		return nil, ErrNilOrderPreviewService
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	auth, err := newJWTPrincipalResolver(jwtSecret, log)
	if err != nil {
		return nil, err
	}

	return &userOrderHandler{
		service: service,
		auth:    auth,
		log:     log.With(logging.String("module", "http-user-order-handler")),
	}, nil
}

func (h *userOrderHandler) RegisterRoutes(router fiber.Router) {
	users := router.Group("/users/:username/orders/:order_id")
	users.Use(h.auth.Middleware())
	users.Get("", h.HandleGetOrderPreviewByID)
}

func (h *userOrderHandler) HandleGetOrderPreviewByID(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)

	username, err := url.PathUnescape(c.Params("username"))
	if err != nil {
		return h.mapError(c, gateway.ErrInvalidUsername)
	}
	username = strings.TrimSpace(username)
	if username == "" {
		return h.mapError(c, gateway.ErrInvalidUsername)
	}
	actorUsername, err := h.auth.Username(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	if strings.TrimSpace(actorUsername) != username {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	userID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	role, err := gateway.ParseParticipantRole(c.Query("role"))
	if err != nil {
		return h.mapError(c, gateway.ErrInvalidParticipantRole)
	}
	req := gateway.GetOrderPreviewByIDRequest{
		OrderID: c.Params("order_id"),
		UserID:  userID,
		Role:    role,
	}
	result, err := h.service.GetOrderPreviewByID(c.UserContext(), req)
	if err != nil {
		return h.mapError(c, err)
	}

	log.Info("user order preview request accepted",
		logging.Operation("http.user.order.preview"),
		logging.DurationMS(time.Since(started)),
		logging.String("username", username),
		logging.String("order_id", req.OrderID),
		logging.String("role", string(role)),
	)
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *userOrderHandler) mapError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, gateway.ErrInvalidUsername),
		errors.Is(err, gateway.ErrInvalidOrderID),
		errors.Is(err, gateway.ErrInvalidParticipantRole),
		errors.Is(err, gateway.ErrInvalidUserID):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrOrderNotOwned):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}
