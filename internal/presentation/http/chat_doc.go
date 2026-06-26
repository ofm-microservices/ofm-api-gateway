package http

import gateway "api-gateway/internal/domain"

var (
	_ = gateway.GetOrderChatResult{}
	_ = gateway.ChatMessage{}
	_ = gateway.ChatAttachment{}
	_ = gateway.CreateChatMessageRequest{}
	_ = gateway.EditChatMessageRequest{}
	_ = gateway.DeleteChatMessageRequest{}
	_ = gateway.CreateChatAttachmentUploadURLRequest{}
	_ = gateway.CreateChatAttachmentUploadURLResult{}
	_ = gateway.CompleteChatAttachmentUploadRequest{}
)

// swaggerGetOrderChatDoc documents the authenticated order chat timeline endpoint.
//
// @Summary Get order chat
// @Description Loads one authenticated participant chat timeline page for an order.
// @Tags chat
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Param order_id path string true "Order ID"
// @Param cursor query string false "Cursor"
// @Param limit query int false "Page size"
// @Success 200 {object} gateway.GetOrderChatResult
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /users/{username}/orders/{order_id}/chat [get]
func swaggerGetOrderChatDoc() {}

// swaggerCreateChatMessageDoc documents the authenticated create-message endpoint.
//
// @Summary Create chat message
// @Description Creates one authenticated order chat message.
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Param order_id path string true "Order ID"
// @Param request body gateway.CreateChatMessageRequest true "Message payload"
// @Success 201 {object} gateway.ChatMessage
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /users/{username}/orders/{order_id}/chat/messages [post]
func swaggerCreateChatMessageDoc() {}

// swaggerEditChatMessageDoc documents the authenticated edit-message endpoint.
//
// @Summary Edit chat message
// @Description Edits one authenticated order chat message owned by the caller.
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Param order_id path string true "Order ID"
// @Param message_id path string true "Message ID"
// @Param request body gateway.EditChatMessageRequest true "Edit payload"
// @Success 200 {object} gateway.ChatMessage
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /users/{username}/orders/{order_id}/chat/messages/{message_id} [patch]
func swaggerEditChatMessageDoc() {}

// swaggerDeleteChatMessageDoc documents the authenticated delete-message endpoint.
//
// @Summary Delete chat message
// @Description Soft-deletes one authenticated order chat message owned by the caller.
// @Tags chat
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Param order_id path string true "Order ID"
// @Param message_id path string true "Message ID"
// @Success 200 {object} gateway.ChatMessage
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /users/{username}/orders/{order_id}/chat/messages/{message_id} [delete]
func swaggerDeleteChatMessageDoc() {}

// swaggerCreateChatAttachmentUploadURLDoc documents the authenticated direct-upload reservation endpoint.
//
// @Summary Create chat attachment upload URL
// @Description Reserves one direct upload for an authenticated order chat attachment.
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Param order_id path string true "Order ID"
// @Param request body gateway.CreateChatAttachmentUploadURLRequest true "Upload reservation payload"
// @Success 201 {object} gateway.CreateChatAttachmentUploadURLResult
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /users/{username}/orders/{order_id}/chat/attachments/upload-url [post]
func swaggerCreateChatAttachmentUploadURLDoc() {}

// swaggerCompleteChatAttachmentUploadDoc documents the authenticated direct-upload completion endpoint.
//
// @Summary Complete chat attachment upload
// @Description Finalizes one authenticated order chat attachment after the client upload succeeds.
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Param order_id path string true "Order ID"
// @Param request body gateway.CompleteChatAttachmentUploadRequest true "Upload completion payload"
// @Success 200 {object} gateway.ChatAttachment
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /users/{username}/orders/{order_id}/chat/attachments/complete [post]
func swaggerCompleteChatAttachmentUploadDoc() {}
