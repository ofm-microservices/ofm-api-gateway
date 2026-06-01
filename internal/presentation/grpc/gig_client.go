package grpc

import (
	gateway "api-gateway/internal/domain"
	"context"

	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"github.com/ofm-microservices/ofm-common/pkg/observability/metrics"
	gigv1 "github.com/ofm-microservices/ofm-common/proto/gig/v1"
	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type gigClient struct {
	conn *grpcpkg.ClientConn
	cl   gigv1.GigCommandServiceClient
	mapr GigMapper
	log  logging.Logger
}

// NewGigClient constructs the gRPC client used by api-gateway to manage gig
// drafts.
func NewGigClient(cfg GigServiceConfig, log Logger) (GigClient, error) {
	if cfg.Address == "" {
		return nil, ErrEmptyGigServiceAddress
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

	return &gigClient{
		conn: conn,
		cl:   gigv1.NewGigCommandServiceClient(conn),
		mapr: newGigMapper(log),
		log:  log.With(logging.String("module", "grpc-gig-client"), logging.String("address", cfg.Address)),
	}, nil
}

func (c *gigClient) CreateDraft(ctx context.Context, req gateway.CreateGigDraftRequest) (*gateway.Gig, error) {
	res, err := c.cl.CreateDraft(ctx, c.mapr.ToCreateDraftRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}

	return c.mapr.ToCreateDraftResponse(res), nil
}

func (c *gigClient) UpdateBasicInfo(ctx context.Context, req gateway.UpdateGigBasicInfoRequest) (*gateway.Gig, error) {
	res, err := c.cl.UpdateBasicInfo(ctx, c.mapr.ToUpdateBasicInfoRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}

	return c.mapr.ToUpdateBasicInfoResponse(res), nil
}

func (c *gigClient) ReplacePackages(ctx context.Context, req gateway.ReplaceGigPackagesRequest) (*gateway.Gig, error) {
	res, err := c.cl.ReplacePackages(ctx, c.mapr.ToReplacePackagesRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}

	return c.mapr.ToReplacePackagesResponse(res), nil
}

func (c *gigClient) ReplaceQuestions(ctx context.Context, req gateway.ReplaceGigQuestionsRequest) (*gateway.Gig, error) {
	res, err := c.cl.ReplaceQuestions(ctx, c.mapr.ToReplaceQuestionsRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}

	return c.mapr.ToReplaceQuestionsResponse(res), nil
}

func (c *gigClient) ReplaceMedia(ctx context.Context, req gateway.ReplaceGigMediaRequest) (*gateway.Gig, error) {
	res, err := c.cl.ReplaceMedia(ctx, c.mapr.ToReplaceMediaRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}

	return c.mapr.ToReplaceMediaResponse(res), nil
}

func (c *gigClient) GetDraft(ctx context.Context, req gateway.GetGigDraftRequest) (*gateway.Gig, error) {
	res, err := c.cl.GetDraft(ctx, c.mapr.ToGetDraftRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}

	return c.mapr.ToGetDraftResponse(res), nil
}

func (c *gigClient) GetBySlug(ctx context.Context, req gateway.GetGigBySlugRequest) (*gateway.Gig, error) {
	res, err := c.cl.GetGigBySlug(ctx, c.mapr.ToGetBySlugRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}

	return c.mapr.ToGetBySlugResponse(res), nil
}

func (c *gigClient) Publish(ctx context.Context, req gateway.PublishGigRequest) (*gateway.Gig, error) {
	res, err := c.cl.Publish(ctx, c.mapr.ToPublishRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}

	return c.mapr.ToPublishResponse(res), nil
}

// Close closes the underlying gig-service gRPC client connection.
func (c *gigClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	c.log.Info("closing gig service grpc client")
	return c.conn.Close()
}
