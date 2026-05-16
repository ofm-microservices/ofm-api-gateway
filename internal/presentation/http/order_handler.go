package http

import (
	gateway "api-gateway/internal/domain"
	"errors"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type orderHandler struct {
	service OrderService
	auth    GigPrincipalResolver
	log     logging.Logger
}

// NewOrderHandler constructs the create-order HTTP handler group.
func NewOrderHandler(service OrderService, jwtSecret string, log logging.Logger) (OrderHandler, error) {
	if service == nil {
		return nil, ErrNilOrderService
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	auth, err := newJWTPrincipalResolver(jwtSecret, log)
	if err != nil {
		return nil, err
	}

	return &orderHandler{
		service: service,
		auth:    auth,
		log:     log.With(logging.String("module", "http-order-handler")),
	}, nil
}

// RegisterRoutes mounts order routes under the router it receives.
func (h *orderHandler) RegisterRoutes(router fiber.Router) {
	orders := router.Group("/orders")
	orders.Use(h.auth.Middleware())
	orders.Post("", h.HandleCreateOrder)
}

// HandleCreateOrder validates and forwards the order start request.
func (h *orderHandler) HandleCreateOrder(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)

	var req gateway.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}

	buyerID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	req.BuyerID = buyerID
	req.RealtimeConnectionID = strings.TrimSpace(firstNonEmpty(c.Get("X-Realtime-Connection-Id"), req.RealtimeConnectionID))

	result, err := h.service.CreateOrder(c.UserContext(), req)
	if err != nil {
		return h.mapOrderError(c, err)
	}

	log.Info("order create request accepted",
		logging.Operation("http.order.create"),
		logging.DurationMS(time.Since(started)),
		logging.String("buyer_id", buyerID),
		logging.String("connection_id", req.RealtimeConnectionID),
	)

	return c.Status(fiber.StatusAccepted).JSON(result)
}

func (h *orderHandler) mapOrderError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, gateway.ErrInvalidOrderConnectionID),
		errors.Is(err, gateway.ErrInvalidOrderBuyerEmail),
		errors.Is(err, gateway.ErrInvalidOrderCurrency),
		errors.Is(err, gateway.ErrInvalidOrderPackage),
		errors.Is(err, gateway.ErrInvalidOrderPrice),
		errors.Is(err, gateway.ErrInvalidOrderGigID),
		errors.Is(err, gateway.ErrInvalidOrderBuyerID),
		errors.Is(err, gateway.ErrInvalidOrderTitle),
		errors.Is(err, gateway.ErrInvalidPackageDeliveryDays):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	default:
		h.log.Error("request failed", logging.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
