package http

import (
	gateway "api-gateway/internal/domain"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
)

type chatHandler struct {
	service ChatService
	auth    GigPrincipalResolver
	log     logging.Logger
}

// NewChatHandler constructs the authenticated user-scoped order chat HTTP handler group.
func NewChatHandler(service ChatService, jwtSecret string, log logging.Logger) (ChatHandler, error) {
	if service == nil {
		return nil, ErrNilChatService
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	auth, err := newJWTPrincipalResolver(jwtSecret, log)
	if err != nil {
		return nil, err
	}

	return &chatHandler{
		service: service,
		auth:    auth,
		log:     log.With(logging.String("module", "http-chat-handler")),
	}, nil
}

func (h *chatHandler) RegisterRoutes(router fiber.Router) {
	chat := router.Group("/users/:username/orders/:order_id/chat")
	chat.Use(h.auth.Middleware())
	chat.Get("", h.HandleGetOrderChat)
	chat.Post("/messages", h.HandleCreateMessage)
	chat.Patch("/messages/:message_id", h.HandleEditMessage)
	chat.Delete("/messages/:message_id", h.HandleDeleteMessage)
	chat.Post("/attachments/upload-url", h.HandleCreateAttachmentUploadURL)
	chat.Post("/attachments/complete", h.HandleCompleteAttachmentUpload)
}

func (h *chatHandler) HandleGetOrderChat(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)

	username, userID, err := h.requireOrderChatPrincipal(c)
	if err != nil {
		return err
	}

	limit := int32(50)
	if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
		parsed, parseErr := strconv.ParseInt(raw, 10, 32)
		if parseErr != nil || parsed < 0 {
			return h.mapError(c, gateway.ErrInvalidChatCursor)
		}
		limit = int32(parsed)
	}

	req := gateway.GetOrderChatRequest{
		OrderID: c.Params("order_id"),
		UserID:  userID,
		Cursor:  strings.TrimSpace(c.Query("cursor")),
		Limit:   limit,
	}
	result, err := h.service.GetOrderChat(c.UserContext(), req)
	if err != nil {
		return h.mapError(c, err)
	}

	log.Info("order chat request accepted",
		logging.Operation("http.chat.get_order_chat"),
		logging.DurationMS(time.Since(started)),
		logging.String("username", username),
		logging.String("order_id", req.OrderID),
	)
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *chatHandler) HandleCreateMessage(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)

	username, userID, err := h.requireOrderChatPrincipal(c)
	if err != nil {
		return err
	}

	var req gateway.CreateChatMessageRequest
	if err := c.BodyParser(&req); err != nil && err.Error() != "EOF" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	req.OrderID = c.Params("order_id")
	req.UserID = userID

	result, err := h.service.CreateMessage(c.UserContext(), req)
	if err != nil {
		return h.mapError(c, err)
	}

	log.Info("chat message create request accepted",
		logging.Operation("http.chat.create_message"),
		logging.DurationMS(time.Since(started)),
		logging.String("username", username),
		logging.String("order_id", req.OrderID),
	)
	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *chatHandler) HandleEditMessage(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)

	username, userID, err := h.requireOrderChatPrincipal(c)
	if err != nil {
		return err
	}

	var req gateway.EditChatMessageRequest
	if err := c.BodyParser(&req); err != nil && err.Error() != "EOF" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	req.OrderID = c.Params("order_id")
	req.UserID = userID
	req.MessageID = c.Params("message_id")

	result, err := h.service.EditMessage(c.UserContext(), req)
	if err != nil {
		return h.mapError(c, err)
	}

	log.Info("chat message edit request accepted",
		logging.Operation("http.chat.edit_message"),
		logging.DurationMS(time.Since(started)),
		logging.String("username", username),
		logging.String("order_id", req.OrderID),
		logging.String("message_id", req.MessageID),
	)
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *chatHandler) HandleDeleteMessage(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)

	username, userID, err := h.requireOrderChatPrincipal(c)
	if err != nil {
		return err
	}

	req := gateway.DeleteChatMessageRequest{
		OrderID:   c.Params("order_id"),
		UserID:    userID,
		MessageID: c.Params("message_id"),
	}
	result, err := h.service.DeleteMessage(c.UserContext(), req)
	if err != nil {
		return h.mapError(c, err)
	}

	log.Info("chat message delete request accepted",
		logging.Operation("http.chat.delete_message"),
		logging.DurationMS(time.Since(started)),
		logging.String("username", username),
		logging.String("order_id", req.OrderID),
		logging.String("message_id", req.MessageID),
	)
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *chatHandler) HandleCreateAttachmentUploadURL(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)

	username, userID, err := h.requireOrderChatPrincipal(c)
	if err != nil {
		return err
	}

	var req gateway.CreateChatAttachmentUploadURLRequest
	if err := c.BodyParser(&req); err != nil && err.Error() != "EOF" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	req.OrderID = c.Params("order_id")
	req.UserID = userID

	result, err := h.service.CreateAttachmentUploadURL(c.UserContext(), req)
	if err != nil {
		return h.mapError(c, err)
	}

	log.Info("chat attachment upload-url request accepted",
		logging.Operation("http.chat.create_attachment_upload_url"),
		logging.DurationMS(time.Since(started)),
		logging.String("username", username),
		logging.String("order_id", req.OrderID),
	)
	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *chatHandler) HandleCompleteAttachmentUpload(c *fiber.Ctx) error {
	started := time.Now()
	log := logging.WithContext(c.UserContext(), h.log)

	username, userID, err := h.requireOrderChatPrincipal(c)
	if err != nil {
		return err
	}

	var req gateway.CompleteChatAttachmentUploadRequest
	if err := c.BodyParser(&req); err != nil && err.Error() != "EOF" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	req.OrderID = c.Params("order_id")
	req.UserID = userID

	result, err := h.service.CompleteAttachmentUpload(c.UserContext(), req)
	if err != nil {
		return h.mapError(c, err)
	}

	log.Info("chat attachment complete request accepted",
		logging.Operation("http.chat.complete_attachment_upload"),
		logging.DurationMS(time.Since(started)),
		logging.String("username", username),
		logging.String("order_id", req.OrderID),
		logging.String("file_id", req.FileID),
	)
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *chatHandler) requireOrderChatPrincipal(c *fiber.Ctx) (string, string, error) {
	username, err := url.PathUnescape(c.Params("username"))
	if err != nil {
		return "", "", h.mapError(c, gateway.ErrInvalidUsername)
	}
	username = strings.TrimSpace(username)
	if username == "" {
		return "", "", h.mapError(c, gateway.ErrInvalidUsername)
	}
	actorUsername, err := h.auth.Username(c)
	if err != nil {
		return "", "", c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	if strings.TrimSpace(actorUsername) != username {
		return "", "", c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	userID, err := h.auth.FreelancerID(c)
	if err != nil {
		return "", "", c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	return username, userID, nil
}

func (h *chatHandler) mapError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, gateway.ErrInvalidUsername),
		errors.Is(err, gateway.ErrInvalidOrderID),
		errors.Is(err, gateway.ErrInvalidUserID),
		errors.Is(err, gateway.ErrInvalidChatMessageID),
		errors.Is(err, gateway.ErrInvalidChatCursor),
		errors.Is(err, gateway.ErrInvalidChatFileID),
		errors.Is(err, gateway.ErrInvalidChatFilename),
		errors.Is(err, gateway.ErrInvalidChatContentType),
		errors.Is(err, gateway.ErrInvalidChatSizeBytes),
		errors.Is(err, gateway.ErrInvalidOrderMessage):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrChatNotFound),
		errors.Is(err, gateway.ErrChatMessageNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrChatAccessDenied):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrChatClosed):
		return c.Status(fiber.StatusPreconditionFailed).JSON(fiber.Map{"error": err.Error()})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}
