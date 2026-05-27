package grpc

import (
	gateway "api-gateway/internal/domain"
	"context"

	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"github.com/ofm-microservices/ofm-common/pkg/observability/metrics"
	searchv1 "github.com/ofm-microservices/ofm-common/proto/search/v1"
	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type searchClient struct {
	conn *grpcpkg.ClientConn
	cl   searchv1.SearchServiceClient
	log  logging.Logger
	mapr SearchMapper
}

// NewSearchClient constructs the gRPC client used by api-gateway to query search-service.
func NewSearchClient(cfg SearchServiceConfig, log Logger) (SearchClient, error) {
	if cfg.Address == "" {
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

	return &searchClient{
		conn: conn,
		cl:   searchv1.NewSearchServiceClient(conn),
		log:  log.With(logging.String("module", "grpc-search-client"), logging.String("address", cfg.Address)),
		mapr: newSearchMapper(),
	}, nil
}

func (c *searchClient) Search(ctx context.Context, req gateway.SearchRequest) (*gateway.SearchResponse, error) {
	res, err := c.cl.Search(ctx, c.mapr.ToSearchRequest(req))
	if err != nil {
		return nil, err
	}
	return c.mapr.ToSearchResponse(res), nil
}

func (c *searchClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	c.log.Info("closing search grpc client")
	return c.conn.Close()
}
