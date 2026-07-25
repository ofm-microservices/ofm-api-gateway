package grpc

import (
	"context"
	"strings"

	gateway "api-gateway/internal/domain"

	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"github.com/ofm-microservices/ofm-common/pkg/observability/metrics"
	ordercheckoutv1 "github.com/ofm-microservices/ofm-common/proto/ordercheckout/v1"
	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type orderCheckoutClient struct {
	conn *grpcpkg.ClientConn
	cl   ordercheckoutv1.OrderCheckoutServiceClient
	log  logging.Logger
	mapr OrderCheckoutMapper
}

// NewOrderCheckoutClient constructs the gRPC client used by api-gateway to
// orchestrate the hybrid order checkout flow.
func NewOrderCheckoutClient(cfg OrderSagaConfig, log Logger) (OrderCheckoutClient, error) {
	if strings.TrimSpace(cfg.Address) == "" {
		return nil, ErrEmptyOrderSagaAddress
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
	return &orderCheckoutClient{
		conn: conn,
		cl:   ordercheckoutv1.NewOrderCheckoutServiceClient(conn),
		log:  log.With(logging.String("module", "order-checkout-client"), logging.String("address", cfg.Address)),
		mapr: newOrderCheckoutMapper(),
	}, nil
}

func (c *orderCheckoutClient) StartOrder(ctx context.Context, req gateway.CreateOrderRequest) (*gateway.CreateOrderResult, error) {
	res, err := c.cl.StartOrder(ctx, c.mapr.ToStartOrderRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToStartOrderResponse(res), nil
}

func (c *orderCheckoutClient) ConfirmOrder(ctx context.Context, req gateway.ConfirmOrderRequest) (*gateway.ConfirmOrderResult, error) {
	res, err := c.cl.ConfirmOrder(ctx, c.mapr.ToConfirmOrderRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToConfirmOrderResponse(res), nil
}

func (c *orderCheckoutClient) SubmitRequirements(ctx context.Context, req gateway.SubmitOrderRequirementsRequest) (*gateway.SubmitOrderRequirementsResult, error) {
	res, err := c.cl.SubmitRequirements(ctx, c.mapr.ToSubmitRequirementsRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToSubmitRequirementsResponse(res), nil
}

func (c *orderCheckoutClient) SubmitMessage(ctx context.Context, req gateway.SubmitOrderMessageRequest) (*gateway.SubmitOrderMessageResult, error) {
	res, err := c.cl.SubmitMessage(ctx, c.mapr.ToSubmitMessageRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToSubmitMessageResponse(res), nil
}

func (c *orderCheckoutClient) CreateAttachmentUploadURL(ctx context.Context, req gateway.CreateOrderAttachmentUploadURLRequest) (*gateway.CreateOrderAttachmentUploadURLResult, error) {
	res, err := c.cl.CreateAttachmentUploadURL(ctx, c.mapr.ToCreateAttachmentUploadURLRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToCreateAttachmentUploadURLResponse(res), nil
}

func (c *orderCheckoutClient) CompleteAttachmentUpload(ctx context.Context, req gateway.CompleteOrderAttachmentUploadRequest) (*gateway.CompleteOrderAttachmentUploadResult, error) {
	res, err := c.cl.CompleteAttachmentUpload(ctx, c.mapr.ToCompleteAttachmentUploadRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToCompleteAttachmentUploadResponse(res), nil
}

func (c *orderCheckoutClient) DeliverOrder(ctx context.Context, req gateway.DeliverOrderRequest) (*gateway.DeliverOrderResult, error) {
	res, err := c.cl.DeliverOrder(ctx, c.mapr.ToDeliverOrderRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToDeliverOrderResponse(res), nil
}

func (c *orderCheckoutClient) AcceptDelivery(ctx context.Context, req gateway.AcceptDeliveryRequest) (*gateway.AcceptDeliveryResult, error) {
	res, err := c.cl.AcceptDelivery(ctx, c.mapr.ToAcceptDeliveryRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToAcceptDeliveryResponse(res), nil
}

func (c *orderCheckoutClient) RequestRevision(ctx context.Context, req gateway.RequestRevisionRequest) (*gateway.RequestRevisionResult, error) {
	res, err := c.cl.RequestRevision(ctx, c.mapr.ToRequestRevisionRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToRequestRevisionResponse(res), nil
}

func (c *orderCheckoutClient) OpenDispute(ctx context.Context, req gateway.OpenDisputeRequest) (*gateway.OpenDisputeResult, error) {
	res, err := c.cl.OpenDispute(ctx, c.mapr.ToOpenDisputeRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToOpenDisputeResponse(res), nil
}

func (c *orderCheckoutClient) ResolveDispute(ctx context.Context, req gateway.ResolveDisputeRequest) (*gateway.ResolveDisputeResult, error) {
	res, err := c.cl.ResolveDispute(ctx, c.mapr.ToResolveDisputeRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToResolveDisputeResponse(res), nil
}

func (c *orderCheckoutClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	c.log.Info("closing order checkout grpc client")
	return c.conn.Close()
}
