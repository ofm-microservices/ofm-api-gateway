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
	orders.Post("/start", h.HandleCreateOrder)
	orders.Post("/:order_id/confirm", h.HandleConfirmOrder)
	orders.Post("/:order_id/requirements", h.HandleSubmitRequirements)
	orders.Post("/:order_id/message", h.HandleSubmitMessage)
	orders.Post("/:order_id/deliver", h.HandleDeliverOrder)
	orders.Post("/:order_id/accept", h.HandleAcceptDelivery)
	orders.Post("/:order_id/request-revision", h.HandleRequestRevision)
	orders.Post("/:order_id/dispute", h.HandleOpenDispute)
	orders.Post("/:order_id/dispute/resolve", h.HandleResolveDispute)
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
	if req.BuyerEmail == "" {
		if buyerEmail, emailErr := h.auth.Email(c); emailErr == nil {
			req.BuyerEmail = buyerEmail
		}
	}

	result, err := h.service.CreateOrder(c.UserContext(), req)
	if err != nil {
		return h.mapOrderError(c, err)
	}

	log.Info("order create request accepted",
		logging.Operation("http.order.create"),
		logging.DurationMS(time.Since(started)),
		logging.String("buyer_id", buyerID),
	)

	return c.Status(fiber.StatusCreated).JSON(result)
}

// HandleConfirmOrder validates and forwards the order confirmation request.
func (h *orderHandler) HandleConfirmOrder(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)

	req := gateway.ConfirmOrderRequest{
		OrderID: c.Params("order_id"),
	}
	buyerID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	if err := c.BodyParser(&req); err != nil && err.Error() != "EOF" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	req.OrderID = c.Params("order_id")
	req.BuyerID = buyerID

	result, err := h.service.ConfirmOrder(c.UserContext(), req)
	if err != nil {
		return h.mapOrderError(c, err)
	}

	log.Info("order confirm request accepted",
		logging.Operation("http.order.confirm"),
		logging.DurationMS(time.Since(started)),
		logging.String("order_id", req.OrderID),
	)

	return c.Status(fiber.StatusOK).JSON(result)
}

// HandleSubmitRequirements validates and forwards the buyer requirements answers request.
func (h *orderHandler) HandleSubmitRequirements(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)

	var req gateway.SubmitOrderRequirementsRequest
	if err := c.BodyParser(&req); err != nil && err.Error() != "EOF" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	buyerID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	req.OrderID = c.Params("order_id")
	req.BuyerID = buyerID

	result, err := h.service.SubmitRequirements(c.UserContext(), req)
	if err != nil {
		return h.mapOrderError(c, err)
	}

	log.Info("order requirements request accepted",
		logging.Operation("http.order.submit_requirements"),
		logging.DurationMS(time.Since(started)),
		logging.String("order_id", req.OrderID),
		logging.String("buyer_id", buyerID),
	)

	return c.Status(fiber.StatusOK).JSON(result)
}

// HandleSubmitMessage validates and forwards the buyer initial message request.
func (h *orderHandler) HandleSubmitMessage(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)

	var req gateway.SubmitOrderMessageRequest
	if err := c.BodyParser(&req); err != nil && err.Error() != "EOF" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	buyerID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	req.OrderID = c.Params("order_id")
	req.BuyerID = buyerID

	result, err := h.service.SubmitMessage(c.UserContext(), req)
	if err != nil {
		return h.mapOrderError(c, err)
	}

	log.Info("order message request accepted",
		logging.Operation("http.order.submit_message"),
		logging.DurationMS(time.Since(started)),
		logging.String("order_id", req.OrderID),
		logging.String("buyer_id", buyerID),
	)

	return c.Status(fiber.StatusOK).JSON(result)
}

// HandleDeliverOrder validates and forwards the seller delivery request.
func (h *orderHandler) HandleDeliverOrder(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)

	var req gateway.DeliverOrderRequest
	if err := c.BodyParser(&req); err != nil && err.Error() != "EOF" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	sellerID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	req.OrderID = c.Params("order_id")
	req.SellerID = sellerID

	result, err := h.service.DeliverOrder(c.UserContext(), req)
	if err != nil {
		return h.mapOrderError(c, err)
	}

	log.Info("order deliver request accepted",
		logging.Operation("http.order.deliver"),
		logging.DurationMS(time.Since(started)),
		logging.String("order_id", req.OrderID),
		logging.String("seller_id", sellerID),
	)

	return c.Status(fiber.StatusOK).JSON(result)
}

// HandleAcceptDelivery validates and forwards the buyer acceptance request.
func (h *orderHandler) HandleAcceptDelivery(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)

	var req gateway.AcceptDeliveryRequest
	if err := c.BodyParser(&req); err != nil && err.Error() != "EOF" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	buyerID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	req.OrderID = c.Params("order_id")
	req.BuyerID = buyerID

	result, err := h.service.AcceptDelivery(c.UserContext(), req)
	if err != nil {
		return h.mapOrderError(c, err)
	}

	log.Info("order accept delivery request accepted",
		logging.Operation("http.order.accept_delivery"),
		logging.DurationMS(time.Since(started)),
		logging.String("order_id", req.OrderID),
		logging.String("buyer_id", buyerID),
	)

	return c.Status(fiber.StatusOK).JSON(result)
}

// HandleRequestRevision validates and forwards the buyer revision request.
func (h *orderHandler) HandleRequestRevision(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)

	var req gateway.RequestRevisionRequest
	if err := c.BodyParser(&req); err != nil && err.Error() != "EOF" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	buyerID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	req.OrderID = c.Params("order_id")
	req.BuyerID = buyerID

	result, err := h.service.RequestRevision(c.UserContext(), req)
	if err != nil {
		return h.mapOrderError(c, err)
	}

	log.Info("order revision request accepted",
		logging.Operation("http.order.request_revision"),
		logging.DurationMS(time.Since(started)),
		logging.String("order_id", req.OrderID),
		logging.String("buyer_id", buyerID),
	)

	return c.Status(fiber.StatusOK).JSON(result)
}

// HandleOpenDispute validates and forwards an owner dispute request.
func (h *orderHandler) HandleOpenDispute(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)

	var req gateway.OpenDisputeRequest
	if err := c.BodyParser(&req); err != nil && err.Error() != "EOF" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	userID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	req.OrderID = c.Params("order_id")
	req.ActorID = userID

	result, err := h.service.OpenDispute(c.UserContext(), req)
	if err != nil {
		return h.mapOrderError(c, err)
	}

	log.Info("order dispute request accepted",
		logging.Operation("http.order.open_dispute"),
		logging.DurationMS(time.Since(started)),
		logging.String("order_id", req.OrderID),
		logging.String("user_id", userID),
	)

	return c.Status(fiber.StatusOK).JSON(result)
}

// HandleResolveDispute validates and forwards the admin dispute settlement request.
func (h *orderHandler) HandleResolveDispute(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)

	var req gateway.ResolveDisputeRequest
	if err := c.BodyParser(&req); err != nil && err.Error() != "EOF" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	adminID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	isAdmin, err := h.auth.HasRole(c, gateway.RoleAdmin)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	req.OrderID = c.Params("order_id")
	req.AdminUserID = adminID

	result, err := h.service.ResolveDispute(c.UserContext(), req)
	if err != nil {
		return h.mapOrderError(c, err)
	}

	log.Info("order dispute resolution accepted",
		logging.Operation("http.order.resolve_dispute"),
		logging.DurationMS(time.Since(started)),
		logging.String("order_id", req.OrderID),
		logging.String("admin_user_id", adminID),
	)

	return c.Status(fiber.StatusOK).JSON(result)
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
	case errors.Is(err, gateway.ErrSelfOrderNotAllowed):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrGigNotFound),
		errors.Is(err, gateway.ErrOrderNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrOrderNotConfirmable),
		errors.Is(err, gateway.ErrOrderRequirementsIncomplete),
		errors.Is(err, gateway.ErrOrderAlreadyPaymentPending),
		errors.Is(err, gateway.ErrOrderAlreadyFunded),
		errors.Is(err, gateway.ErrOrderNotDeliverable),
		errors.Is(err, gateway.ErrOrderNotAcceptable),
		errors.Is(err, gateway.ErrOrderNotRevisionable),
		errors.Is(err, gateway.ErrOrderNotDisputable):
		return c.Status(fiber.StatusPreconditionFailed).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrInvalidDisputeSplit):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrOrderReleaseFailed):
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrOrderNotOwned):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrInvalidOrderDeliveryMessage),
		errors.Is(err, gateway.ErrInvalidOrderMessage),
		errors.Is(err, gateway.ErrInvalidOrderReason),
		errors.Is(err, gateway.ErrInvalidOrderID),
		errors.Is(err, gateway.ErrInvalidOrderSellerID):
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
