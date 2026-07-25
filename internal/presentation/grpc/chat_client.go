package grpc

import (
	"context"
	"strings"

	gateway "api-gateway/internal/domain"

	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"github.com/ofm-microservices/ofm-common/pkg/observability/metrics"
	chatv1 "github.com/ofm-microservices/ofm-common/proto/chat/v1"
	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type chatClient struct {
	conn *grpcpkg.ClientConn
	cl   chatv1.ChatServiceClient
	log  logging.Logger
	mapr ChatMapper
}

// NewChatClient constructs the gRPC client used by api-gateway for order chats.
func NewChatClient(cfg ChatServiceConfig, log Logger) (ChatClient, error) {
	if strings.TrimSpace(cfg.Address) == "" {
		return nil, ErrEmptyChatAddress
	}
	if log == nil {
		return nil, ErrNilLogger
	}
	conn, err := grpcpkg.NewClient(
		cfg.Address,
		grpcpkg.WithTransportCredentials(insecure.NewCredentials()),
		grpcpkg.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpcpkg.WithUnaryInterceptor(metrics.UnaryClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}
	return &chatClient{
		conn: conn,
		cl:   chatv1.NewChatServiceClient(conn),
		log:  log.With(logging.String("module", "chat-client"), logging.String("address", cfg.Address)),
		mapr: newChatMapper(),
	}, nil
}

func (c *chatClient) GetOrderChat(ctx context.Context, req gateway.GetOrderChatRequest) (*gateway.GetOrderChatResult, error) {
	res, err := c.cl.GetOrderChat(ctx, c.mapr.ToGetOrderChatRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToGetOrderChatResponse(res), nil
}

func (c *chatClient) CreateMessage(ctx context.Context, req gateway.CreateChatMessageRequest) (*gateway.ChatMessage, error) {
	res, err := c.cl.CreateMessage(ctx, c.mapr.ToCreateMessageRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToCreateMessageResponse(res), nil
}

func (c *chatClient) EditMessage(ctx context.Context, req gateway.EditChatMessageRequest) (*gateway.ChatMessage, error) {
	res, err := c.cl.EditMessage(ctx, c.mapr.ToEditMessageRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToEditMessageResponse(res), nil
}

func (c *chatClient) DeleteMessage(ctx context.Context, req gateway.DeleteChatMessageRequest) (*gateway.ChatMessage, error) {
	res, err := c.cl.DeleteMessage(ctx, c.mapr.ToDeleteMessageRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToDeleteMessageResponse(res), nil
}

func (c *chatClient) CreateAttachmentUploadURL(ctx context.Context, req gateway.CreateChatAttachmentUploadURLRequest) (*gateway.CreateChatAttachmentUploadURLResult, error) {
	res, err := c.cl.CreateAttachmentUploadURL(ctx, c.mapr.ToCreateAttachmentUploadURLRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToCreateAttachmentUploadURLResponse(res), nil
}

func (c *chatClient) CompleteAttachmentUpload(ctx context.Context, req gateway.CompleteChatAttachmentUploadRequest) (*gateway.ChatAttachment, error) {
	res, err := c.cl.CompleteAttachmentUpload(ctx, c.mapr.ToCompleteAttachmentUploadRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToCompleteAttachmentUploadResponse(res), nil
}

func (c *chatClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	c.log.Info("closing chat grpc client")
	return c.conn.Close()
}
