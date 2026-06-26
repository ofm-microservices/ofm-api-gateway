package gateway

// GetOrderChatRequest loads a participant-scoped order chat timeline.
type GetOrderChatRequest struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
	Cursor  string `json:"cursor"`
	Limit   int32  `json:"limit"`
}

// CreateChatMessageRequest stores one new order chat message.
type CreateChatMessageRequest struct {
	OrderID       string   `json:"order_id"`
	UserID        string   `json:"user_id"`
	Text          string   `json:"text"`
	AttachmentIDs []string `json:"attachment_ids"`
}

// EditChatMessageRequest updates one existing order chat message.
type EditChatMessageRequest struct {
	OrderID   string `json:"order_id"`
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
	Text      string `json:"text"`
}

// DeleteChatMessageRequest soft-deletes one order chat message.
type DeleteChatMessageRequest struct {
	OrderID   string `json:"order_id"`
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
}

// CreateChatAttachmentUploadURLRequest reserves one direct upload for a chat attachment.
type CreateChatAttachmentUploadURLRequest struct {
	OrderID     string `json:"order_id"`
	UserID      string `json:"user_id"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	SizeBytes   int64  `json:"size_bytes"`
}

// CompleteChatAttachmentUploadRequest finalizes one direct upload after the client PUT succeeds.
type CompleteChatAttachmentUploadRequest struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
	FileID  string `json:"file_id"`
}

// ChatAttachment describes one chat attachment returned to the frontend.
type ChatAttachment struct {
	FileID      string `json:"file_id"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	SizeBytes   int64  `json:"size_bytes"`
	URL         string `json:"url"`
}

// ChatMessage describes one order chat message returned to the frontend.
type ChatMessage struct {
	MessageID    string           `json:"message_id"`
	OrderID      string           `json:"order_id"`
	SenderUserID string           `json:"sender_user_id"`
	MessageType  string           `json:"message_type"`
	Text         string           `json:"text"`
	Deleted      bool             `json:"deleted"`
	Edited       bool             `json:"edited"`
	Attachments  []ChatAttachment `json:"attachments"`
	CreatedAt    string           `json:"created_at"`
	UpdatedAt    string           `json:"updated_at"`
}

// GetOrderChatResult returns one chat timeline page.
type GetOrderChatResult struct {
	OrderID     string        `json:"order_id"`
	BuyerID     string        `json:"buyer_id"`
	SellerID    string        `json:"seller_id"`
	Status      string        `json:"status"`
	CloseReason string        `json:"close_reason"`
	Messages    []ChatMessage `json:"messages"`
	NextCursor  string        `json:"next_cursor"`
}

// CreateChatAttachmentUploadURLResult returns one direct upload reservation.
type CreateChatAttachmentUploadURLResult struct {
	FileID    string            `json:"file_id"`
	UploadURL string            `json:"upload_url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	ExpiresAt string            `json:"expires_at"`
}
