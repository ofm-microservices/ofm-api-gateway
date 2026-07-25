package grpc

import (
	"strings"

	gateway "api-gateway/internal/domain"

	chatv1 "github.com/ofm-microservices/ofm-common/proto/chat/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type chatMapper struct{}

func newChatMapper() ChatMapper { return &chatMapper{} }

func (m *chatMapper) ToGetOrderChatRequest(req gateway.GetOrderChatRequest) *chatv1.GetOrderChatRequest {
	return &chatv1.GetOrderChatRequest{
		OrderId: strings.TrimSpace(req.OrderID),
		UserId:  strings.TrimSpace(req.UserID),
		Cursor:  strings.TrimSpace(req.Cursor),
		Limit:   req.Limit,
	}
}

func (m *chatMapper) ToGetOrderChatResponse(res *chatv1.GetOrderChatResponse) *gateway.GetOrderChatResult {
	if res == nil {
		return &gateway.GetOrderChatResult{}
	}
	msgs := make([]gateway.ChatMessage, 0, len(res.GetMessages()))
	for _, item := range res.GetMessages() {
		msgs = append(msgs, gateway.ChatMessage{
			MessageID:    item.GetMessageId(),
			OrderID:      item.GetOrderId(),
			SenderUserID: item.GetSenderUserId(),
			MessageType:  item.GetMessageType(),
			Text:         item.GetText(),
			Deleted:      item.GetDeleted(),
			Edited:       item.GetEdited(),
			Attachments:  toChatAttachments(item.GetAttachments()),
			CreatedAt:    item.GetCreatedAt(),
			UpdatedAt:    item.GetUpdatedAt(),
		})
	}
	return &gateway.GetOrderChatResult{
		OrderID:     res.GetOrderId(),
		BuyerID:     res.GetBuyerId(),
		SellerID:    res.GetSellerId(),
		Status:      res.GetStatus(),
		CloseReason: res.GetCloseReason(),
		Messages:    msgs,
		NextCursor:  res.GetNextCursor(),
	}
}

func (m *chatMapper) ToCreateMessageRequest(req gateway.CreateChatMessageRequest) *chatv1.CreateMessageRequest {
	return &chatv1.CreateMessageRequest{
		OrderId:       strings.TrimSpace(req.OrderID),
		UserId:        strings.TrimSpace(req.UserID),
		Text:          strings.TrimSpace(req.Text),
		AttachmentIds: append([]string(nil), req.AttachmentIDs...),
	}
}

func (m *chatMapper) ToCreateMessageResponse(res *chatv1.CreateMessageResponse) *gateway.ChatMessage {
	if res == nil || res.GetMessage() == nil {
		return &gateway.ChatMessage{}
	}
	msg := res.GetMessage()
	return &gateway.ChatMessage{
		MessageID:    msg.GetMessageId(),
		OrderID:      msg.GetOrderId(),
		SenderUserID: msg.GetSenderUserId(),
		MessageType:  msg.GetMessageType(),
		Text:         msg.GetText(),
		Deleted:      msg.GetDeleted(),
		Edited:       msg.GetEdited(),
		Attachments:  toChatAttachments(msg.GetAttachments()),
		CreatedAt:    msg.GetCreatedAt(),
		UpdatedAt:    msg.GetUpdatedAt(),
	}
}

func (m *chatMapper) ToEditMessageRequest(req gateway.EditChatMessageRequest) *chatv1.EditMessageRequest {
	return &chatv1.EditMessageRequest{
		OrderId:   strings.TrimSpace(req.OrderID),
		UserId:    strings.TrimSpace(req.UserID),
		MessageId: strings.TrimSpace(req.MessageID),
		Text:      strings.TrimSpace(req.Text),
	}
}

func (m *chatMapper) ToEditMessageResponse(res *chatv1.EditMessageResponse) *gateway.ChatMessage {
	return m.ToCreateMessageResponse(&chatv1.CreateMessageResponse{Message: res.GetMessage()})
}

func (m *chatMapper) ToDeleteMessageRequest(req gateway.DeleteChatMessageRequest) *chatv1.DeleteMessageRequest {
	return &chatv1.DeleteMessageRequest{
		OrderId:   strings.TrimSpace(req.OrderID),
		UserId:    strings.TrimSpace(req.UserID),
		MessageId: strings.TrimSpace(req.MessageID),
	}
}

func (m *chatMapper) ToDeleteMessageResponse(res *chatv1.DeleteMessageResponse) *gateway.ChatMessage {
	return m.ToCreateMessageResponse(&chatv1.CreateMessageResponse{Message: res.GetMessage()})
}

func (m *chatMapper) ToCreateAttachmentUploadURLRequest(req gateway.CreateChatAttachmentUploadURLRequest) *chatv1.CreateAttachmentUploadURLRequest {
	return &chatv1.CreateAttachmentUploadURLRequest{
		OrderId:     strings.TrimSpace(req.OrderID),
		UserId:      strings.TrimSpace(req.UserID),
		Filename:    strings.TrimSpace(req.Filename),
		ContentType: strings.TrimSpace(req.ContentType),
		SizeBytes:   req.SizeBytes,
	}
}

func (m *chatMapper) ToCreateAttachmentUploadURLResponse(res *chatv1.CreateAttachmentUploadURLResponse) *gateway.CreateChatAttachmentUploadURLResult {
	if res == nil {
		return &gateway.CreateChatAttachmentUploadURLResult{}
	}
	return &gateway.CreateChatAttachmentUploadURLResult{
		FileID:    res.GetFileId(),
		UploadURL: res.GetUploadUrl(),
		Method:    res.GetMethod(),
		Headers:   res.GetHeaders(),
		ExpiresAt: res.GetExpiresAt(),
	}
}

func (m *chatMapper) ToCompleteAttachmentUploadRequest(req gateway.CompleteChatAttachmentUploadRequest) *chatv1.CompleteAttachmentUploadRequest {
	return &chatv1.CompleteAttachmentUploadRequest{
		OrderId: strings.TrimSpace(req.OrderID),
		UserId:  strings.TrimSpace(req.UserID),
		FileId:  strings.TrimSpace(req.FileID),
	}
}

func (m *chatMapper) ToCompleteAttachmentUploadResponse(res *chatv1.CompleteAttachmentUploadResponse) *gateway.ChatAttachment {
	if res == nil || res.GetAttachment() == nil {
		return &gateway.ChatAttachment{}
	}
	item := res.GetAttachment()
	return &gateway.ChatAttachment{
		FileID:      item.GetFileId(),
		Filename:    item.GetFilename(),
		ContentType: item.GetContentType(),
		SizeBytes:   item.GetSizeBytes(),
		URL:         item.GetUrl(),
	}
}

func (m *chatMapper) ToError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	switch st.Code() {
	case codes.InvalidArgument:
		switch st.Message() {
		case gateway.ErrInvalidChatCursor.Error():
			return gateway.ErrInvalidChatCursor
		case gateway.ErrInvalidOrderMessage.Error():
			return gateway.ErrInvalidOrderMessage
		case gateway.ErrInvalidOrderID.Error():
			return gateway.ErrInvalidOrderID
		case gateway.ErrInvalidUserID.Error():
			return gateway.ErrInvalidUserID
		default:
			return gateway.ErrInvalidOrderMessage
		}
	case codes.NotFound:
		switch st.Message() {
		case gateway.ErrChatNotFound.Error():
			return gateway.ErrChatNotFound
		case gateway.ErrChatMessageNotFound.Error():
			return gateway.ErrChatMessageNotFound
		default:
			return err
		}
	case codes.PermissionDenied:
		return gateway.ErrChatAccessDenied
	case codes.FailedPrecondition:
		return gateway.ErrChatClosed
	default:
		return err
	}
}

func toChatAttachments(items []*chatv1.ChatAttachment) []gateway.ChatAttachment {
	out := make([]gateway.ChatAttachment, 0, len(items))
	for _, item := range items {
		out = append(out, gateway.ChatAttachment{
			FileID:      item.GetFileId(),
			Filename:    item.GetFilename(),
			ContentType: item.GetContentType(),
			SizeBytes:   item.GetSizeBytes(),
			URL:         item.GetUrl(),
		})
	}
	return out
}
