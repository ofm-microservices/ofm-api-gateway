package grpc

import (
	"context"
	"strings"

	gateway "api-gateway/internal/domain"

	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"github.com/ofm-microservices/ofm-common/pkg/observability/metrics"
	paymentcheckoutv1 "github.com/ofm-microservices/ofm-common/proto/paymentcheckout/v1"
	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type paymentByOrderClient struct {
	conn *grpcpkg.ClientConn
	cl   paymentcheckoutv1.PaymentCheckoutServiceClient
	log  logging.Logger
}

// NewPaymentByOrderClient constructs the gRPC client used by api-gateway to load payment snapshots by order.
func NewPaymentByOrderClient(cfg PaymentServiceConfig, log Logger) (PaymentByOrderClient, error) {
	if strings.TrimSpace(cfg.Address) == "" {
		return nil, ErrEmptyAddress
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
	return &paymentByOrderClient{
		conn: conn,
		cl:   paymentcheckoutv1.NewPaymentCheckoutServiceClient(conn),
		log:  log.With(logging.String("module", "payment-by-order-client"), logging.String("address", cfg.Address)),
	}, nil
}

func (c *paymentByOrderClient) GetPaymentByOrderId(ctx context.Context, orderID string) (*gateway.OrderPreviewPayment, error) {
	res, err := c.cl.GetPaymentByOrderId(ctx, &paymentcheckoutv1.GetPaymentByOrderIdRequest{OrderId: strings.TrimSpace(orderID)})
	if err != nil {
		return nil, err
	}
	return &gateway.OrderPreviewPayment{
		PaymentID:   res.GetPaymentId(),
		AmountCents: res.GetAmountCents(),
		Currency:    res.GetCurrency(),
		CreatedAt:   res.GetCreatedAt(),
		UpdatedAt:   res.GetUpdatedAt(),
	}, nil
}

func (c *paymentByOrderClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	c.log.Info("closing payment-by-order grpc client")
	return c.conn.Close()
}
