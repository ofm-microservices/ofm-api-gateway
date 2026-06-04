package grpc

import (
	"context"
	"strings"

	gateway "api-gateway/internal/domain"

	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"github.com/ofm-microservices/ofm-common/pkg/observability/metrics"
	orderwritev1 "github.com/ofm-microservices/ofm-common/proto/orderwrite/v1"
	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type orderPreviewClient struct {
	conn *grpcpkg.ClientConn
	cl   orderwritev1.OrderWriteServiceClient
	log  logging.Logger
	mapr OrderPreviewMapper
}

// NewOrderPreviewClient constructs the gRPC client used by api-gateway to load a user-scoped order preview.
func NewOrderPreviewClient(cfg OrderServiceConfig, log Logger) (OrderPreviewClient, error) {
	if strings.TrimSpace(cfg.Address) == "" {
		return nil, ErrEmptyOrderAddress
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
	return &orderPreviewClient{
		conn: conn,
		cl:   orderwritev1.NewOrderWriteServiceClient(conn),
		log:  log.With(logging.String("module", "order-preview-client"), logging.String("address", cfg.Address)),
		mapr: newOrderPreviewMapper(),
	}, nil
}

func (c *orderPreviewClient) GetOrderPreviewByID(ctx context.Context, req gateway.GetOrderPreviewByIDRequest) (*gateway.GetOrderPreviewByIDResult, error) {
	res, err := c.cl.GetOrderPreviewByID(ctx, c.mapr.ToGetOrderPreviewByIDRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToGetOrderPreviewByIDResponse(res), nil
}

func (c *orderPreviewClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
