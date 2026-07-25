package grpc

import (
	gateway "api-gateway/internal/domain"
	"context"

	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"github.com/ofm-microservices/ofm-common/pkg/observability/metrics"
	paymentconnectv1 "github.com/ofm-microservices/ofm-common/proto/paymentconnect/v1"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type paymentClient struct {
	conn *grpcpkg.ClientConn
	cl   paymentconnectv1.PaymentOnboardingServiceClient
	mapr *paymentMapper
	log  logging.Logger
}

// NewPaymentOnboardingClient constructs the gRPC client used by api-gateway
// to start freelancer Stripe onboarding.
func NewPaymentOnboardingClient(cfg PaymentServiceConfig, log Logger) (PaymentOnboardingClient, error) {
	if cfg.Address == "" {
		return nil, ErrEmptyAddress
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	conn, err := grpcpkg.NewClient(
		cfg.Address,
		grpcpkg.WithTransportCredentials(insecure.NewCredentials()),
		grpcpkg.WithUnaryInterceptor(metrics.UnaryClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}
	return &paymentClient{
		conn: conn,
		cl:   paymentconnectv1.NewPaymentOnboardingServiceClient(conn),
		mapr: newPaymentMapper(),
		log:  log.With(logging.String("module", "grpc-payment-onboarding-client"), logging.String("address", cfg.Address)),
	}, nil
}

func (c *paymentClient) StartFreelancerOnboarding(ctx context.Context, req gateway.StartFreelancerOnboardingRequest) (*gateway.StartFreelancerOnboardingResult, error) {
	res, err := c.cl.StartFreelancerOnboarding(ctx, c.mapr.ToStartFreelancerOnboardingRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}
	return c.mapr.ToStartFreelancerOnboardingResponse(res), nil
}

func (c *paymentClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	c.log.Info("closing payment onboarding grpc client")
	return c.conn.Close()
}
