package grpc

import (
	gateway "api-gateway/internal/domain"
	"context"

	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"github.com/ofm-microservices/ofm-common/pkg/observability/metrics"
	authv1 "github.com/ofm-microservices/ofm-common/proto/auth/v1"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authSessionClient struct {
	conn *grpcpkg.ClientConn
	cl   authv1.AuthSessionServiceClient
	mapr AuthSessionMapper
	log  logging.Logger
}

// NewAuthSessionClient constructs the gRPC client used by api-gateway to sign
// in against auth-service.
func NewAuthSessionClient(cfg AuthServiceConfig, log Logger) (AuthSessionClient, error) {
	if cfg.Address == "" {
		return nil, ErrEmptyAddress
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	conn, err := grpcpkg.NewClient(cfg.Address, grpcpkg.WithTransportCredentials(insecure.NewCredentials()), grpcpkg.WithUnaryInterceptor(metrics.UnaryClientInterceptor()))
	if err != nil {
		return nil, err
	}

	return &authSessionClient{
		conn: conn,
		cl:   authv1.NewAuthSessionServiceClient(conn),
		mapr: newAuthSessionMapper(),
		log:  log.With(logging.String("module", "grpc-auth-session-client"), logging.String("address", cfg.Address)),
	}, nil
}

// SignIn asks auth-service to validate credentials and return auth-owned tokens.
func (c *authSessionClient) SignIn(ctx context.Context, req gateway.SignInRequest) (*AuthTokensResult, error) {
	response, err := c.cl.SignIn(ctx, c.mapr.ToSignInRequest(req))
	if err != nil {
		return nil, c.mapr.ToError(err)
	}

	return c.mapr.ToSignInResult(response), nil
}

// Refresh asks auth-service to rotate a refresh token and return new tokens.
func (c *authSessionClient) Refresh(ctx context.Context, req gateway.RefreshTokensRequest) (*AuthTokensResult, error) {
	response, err := c.cl.Refresh(ctx, c.mapr.ToRefreshRequest(req))
	if err != nil {
		return nil, c.mapr.ToRefreshError(err)
	}

	return c.mapr.ToRefreshResult(response), nil
}

// SignOut asks auth-service to revoke a refresh token and end the session.
func (c *authSessionClient) SignOut(ctx context.Context, req gateway.SignOutRequest) error {
	_, err := c.cl.SignOut(ctx, c.mapr.ToSignOutRequest(req))
	if err != nil {
		return c.mapr.ToSignOutError(err)
	}

	return nil
}

// Close closes the underlying auth-session gRPC client connection.
func (c *authSessionClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	c.log.Info("closing auth session grpc client")
	return c.conn.Close()
}
