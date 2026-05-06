package grpc

import (
	gateway "api-gateway/internal/domain"
	"context"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"
	registrationv1 "github.com/ofm-microseervices/ofm-common/proto/registration/v1"

	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	conn *grpcpkg.ClientConn
	cl   registrationv1.RegistrationServiceClient
	mapr RegistrationMapper
	log  logging.Logger
}

// NewClient constructs the gRPC client used by api-gateway to start
// registration sessions.
func NewClient(cfg RegistrationSagaConfig, log Logger) (Client, error) {
	if cfg.Address == "" {
		return nil, ErrEmptyAddress
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	conn, err := grpcpkg.NewClient(
		cfg.Address,
		grpcpkg.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &client{
		conn: conn,
		cl:   registrationv1.NewRegistrationServiceClient(conn),
		mapr: newRegistrationMapper(),
		log:  log.With(logging.String("module", "grpc-registration-client"), logging.String("address", cfg.Address)),
	}, nil
}

// StartRegistration forwards the gateway signup payload to the registration
// saga and maps transport errors back to gateway-domain errors.
func (c *client) StartRegistration(ctx context.Context, req SignUpRequest) (*SignUpResult, error) {
	response, err := c.cl.StartRegistration(ctx, c.mapr.ToStartRegistrationRequest(req))
	if err != nil {
		return nil, c.mapr.ToStartRegistrationError(err)
	}

	return c.mapr.ToSignUpResult(response), nil
}

// VerifyEmail forwards an email verification command to the registration saga.
func (c *client) VerifyEmail(ctx context.Context, req VerifyEmailRequest) (*VerifyEmailResult, error) {
	response, err := c.cl.VerifyEmail(ctx, c.mapr.ToVerifyEmailRequest(req))
	if err != nil {
		return nil, c.mapr.ToStartRegistrationError(err)
	}

	return c.mapr.ToVerifyEmailResult(response), nil
}

// GetRegistrationStatus reads the saga state before token completion.
func (c *client) GetRegistrationStatus(ctx context.Context, sessionID, clientID string) (*gateway.RegistrationStatus, error) {
	response, err := c.cl.GetRegistrationStatus(ctx, &registrationv1.GetRegistrationStatusRequest{
		SessionId: sessionID,
		ClientId:  clientID,
	})
	if err != nil {
		return nil, c.mapr.ToRegistrationStatusError(err)
	}

	return c.mapr.ToRegistrationStatus(response), nil
}

// Close closes the underlying gRPC client connection.
func (c *client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	c.log.Info("closing registration saga grpc client")
	return c.conn.Close()
}
