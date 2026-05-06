package grpc

import (
	"context"

	"github.com/ofm-microseervices/ofm-common/pkg/logging"
	authv1 "github.com/ofm-microseervices/ofm-common/proto/auth/v1"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authClient struct {
	conn *grpcpkg.ClientConn
	cl   authv1.AuthQueryServiceClient
	mapr RegistrationMapper
	log  logging.Logger
}

// NewAuthClient constructs the gRPC client used by api-gateway to exchange a
// completed registration for login tokens.
func NewAuthClient(cfg AuthServiceConfig, log Logger) (AuthClient, error) {
	if cfg.Address == "" {
		return nil, ErrEmptyAddress
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	conn, err := grpcpkg.NewClient(cfg.Address, grpcpkg.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &authClient{
		conn: conn,
		cl:   authv1.NewAuthQueryServiceClient(conn),
		mapr: newRegistrationMapper(),
		log:  log.With(logging.String("module", "grpc-auth-client"), logging.String("address", cfg.Address)),
	}, nil
}

// IssueRegistrationTokens asks auth-service to issue auth-owned tokens.
func (c *authClient) IssueRegistrationTokens(ctx context.Context, userID string) (*CompleteRegistrationResult, error) {
	response, err := c.cl.IssueRegistrationTokens(ctx, &authv1.IssueRegistrationTokensRequest{UserId: userID})
	if err != nil {
		return nil, err
	}

	return c.mapr.ToCompleteRegistrationResult(response), nil
}

// Close closes the underlying auth-service gRPC client connection.
func (c *authClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	c.log.Info("closing auth service grpc client")
	return c.conn.Close()
}
