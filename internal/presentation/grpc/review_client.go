package grpc

import (
	"context"
	"strings"

	gateway "api-gateway/internal/domain"

	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"github.com/ofm-microservices/ofm-common/pkg/observability/metrics"
	reviewv1 "github.com/ofm-microservices/ofm-common/proto/review/v1"
	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type reviewClient struct {
	conn *grpcpkg.ClientConn
	cl   reviewv1.ReviewServiceClient
	log  logging.Logger
	mapr ReviewMapper
}

// NewReviewClient constructs the gRPC client used by api-gateway to submit reviews.
func NewReviewClient(cfg ReviewServiceConfig, log Logger) (ReviewClient, error) {
	if strings.TrimSpace(cfg.Address) == "" {
		return nil, ErrEmptyAddress
	}
	if log == nil {
		return nil, ErrNilLogger
	}
	conn, err := grpcpkg.NewClient(cfg.Address, grpcpkg.WithTransportCredentials(insecure.NewCredentials()), grpcpkg.WithStatsHandler(otelgrpc.NewClientHandler()), grpcpkg.WithUnaryInterceptor(metrics.UnaryClientInterceptor()))
	if err != nil {
		return nil, err
	}
	return &reviewClient{
		conn: conn,
		cl:   reviewv1.NewReviewServiceClient(conn),
		log:  log.With(logging.String("module", "review-client"), logging.String("address", cfg.Address)),
		mapr: newReviewMapper(),
	}, nil
}

func (c *reviewClient) CreateReview(ctx context.Context, req gateway.CreateReviewRequest) (*gateway.CreateReviewResult, error) {
	res, err := c.cl.CreateReview(ctx, c.mapr.ToCreateReviewRequest(req, req.BuyerID))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToCreateReviewResponse(res), nil
}

func (c *reviewClient) GetGigReviews(ctx context.Context, req gateway.GetGigReviewsRequest) (*gateway.GetGigReviewsResult, error) {
	res, err := c.cl.ListGigReviews(ctx, c.mapr.ToGetGigReviewsRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToGetGigReviewsResponse(res), nil
}

func (c *reviewClient) GetGigReviewsSummary(ctx context.Context, req gateway.GetGigReviewsSummaryRequest) (*gateway.ReviewSummary, error) {
	res, err := c.cl.GetGigRatingSummary(ctx, c.mapr.ToGetGigReviewsSummaryRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToGetGigReviewsSummaryResponse(res), nil
}

func (c *reviewClient) GetUserRatingSummaryByUsername(ctx context.Context, req gateway.GetUserRatingSummaryByUsernameRequest) (*gateway.ReviewSummary, error) {
	res, err := c.cl.GetUserRatingSummaryByUsername(ctx, c.mapr.ToGetUserRatingSummaryByUsernameRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToGetUserRatingSummaryByUsernameResponse(res), nil
}

func (c *reviewClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
