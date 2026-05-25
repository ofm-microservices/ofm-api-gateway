package http

import (
	gateway "api-gateway/internal/domain"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
)

type reviewHandler struct {
	service ReviewService
	auth    GigPrincipalResolver
	log     logging.Logger
}

// NewReviewHandler constructs the buyer review HTTP handler group.
func NewReviewHandler(service ReviewService, jwtSecret string, log logging.Logger) (ReviewHandler, error) {
	if service == nil {
		return nil, ErrNilReviewService
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	auth, err := newJWTPrincipalResolver(jwtSecret, log)
	if err != nil {
		return nil, err
	}

	return &reviewHandler{
		service: service,
		auth:    auth,
		log:     log.With(logging.String("module", "http-review-handler")),
	}, nil
}

func (h *reviewHandler) RegisterRoutes(router fiber.Router) {
	reviews := router.Group("")
	reviews.Use(h.auth.Middleware())
	reviews.Post("/orders/:order_id/reviews", h.HandleCreateReview)
}

func (h *reviewHandler) HandleCreateReview(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)

	var req gateway.CreateReviewRequest
	if err := c.BodyParser(&req); err != nil && err.Error() != "EOF" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}

	buyerID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	req.OrderID = c.Params("order_id")
	req.BuyerID = buyerID
	req.RequestedAt = strings.TrimSpace(firstNonEmpty(c.Get("X-Requested-At"), req.RequestedAt))

	result, err := h.service.CreateReview(c.UserContext(), req)
	if err != nil {
		return h.mapReviewError(c, err)
	}

	log.Info("order review request accepted",
		logging.Operation("http.review.create"),
		logging.DurationMS(time.Since(started)),
		logging.String("order_id", req.OrderID),
		logging.String("buyer_id", buyerID),
	)

	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *reviewHandler) mapReviewError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, gateway.ErrInvalidOrderID),
		errors.Is(err, gateway.ErrInvalidReviewContent):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrOrderNotOwned):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrReviewOwnerMismatch):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrOrderNotAcceptable):
		return c.Status(fiber.StatusPreconditionFailed).JSON(fiber.Map{"error": err.Error()})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}
