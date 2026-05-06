package grpc

import (
	"api-gateway/config"
	gateway "api-gateway/internal/domain"
	"context"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"
	authv1 "github.com/ofm-microseervices/ofm-common/proto/auth/v1"
	registrationv1 "github.com/ofm-microseervices/ofm-common/proto/registration/v1"
)

// Logger aliases the shared logger contract used by the gRPC adapter.
type Logger = logging.Logger

// SignUpRequest aliases the gateway-domain signup request transported over gRPC.
type SignUpRequest = gateway.SignUpRequest

// SignUpResult aliases the gateway-domain signup result returned from the saga.
type SignUpResult = gateway.SignUpResult

// VerifyEmailRequest aliases the gateway-domain verification command.
type VerifyEmailRequest = gateway.VerifyEmailRequest

// VerifyEmailResult aliases the gateway-domain verification result.
type VerifyEmailResult = gateway.VerifyEmailResult

// CompleteRegistrationResult aliases the gateway-domain token result.
type CompleteRegistrationResult = gateway.CompleteRegistrationResult

// Client is the gateway-facing gRPC adapter for registration-saga-service.
type Client interface {
	StartRegistration(ctx context.Context, req SignUpRequest) (*SignUpResult, error)
	VerifyEmail(ctx context.Context, req VerifyEmailRequest) (*VerifyEmailResult, error)
	GetRegistrationStatus(ctx context.Context, sessionID, clientID string) (*gateway.RegistrationStatus, error)
	Close() error
}

// AuthClient is the gateway-facing gRPC adapter for auth-service token issuing.
type AuthClient interface {
	IssueRegistrationTokens(ctx context.Context, userID string) (*CompleteRegistrationResult, error)
	Close() error
}

// RegistrationMapper translates between gateway-domain signup types and the
// shared registration gRPC contract.
type RegistrationMapper interface {
	ToStartRegistrationRequest(req SignUpRequest) *registrationv1.StartRegistrationRequest
	ToSignUpResult(res *registrationv1.StartRegistrationResponse) *SignUpResult
	ToVerifyEmailRequest(req VerifyEmailRequest) *registrationv1.VerifyEmailRequest
	ToVerifyEmailResult(res *registrationv1.VerifyEmailResponse) *VerifyEmailResult
	ToRegistrationStatus(res *registrationv1.GetRegistrationStatusResponse) *gateway.RegistrationStatus
	ToCompleteRegistrationResult(res *authv1.IssueRegistrationTokensResponse) *CompleteRegistrationResult
	ToStartRegistrationError(err error) error
	ToRegistrationStatusError(err error) error
}

// RegistrationSagaConfig aliases the outbound saga gRPC client configuration.
type RegistrationSagaConfig = config.RegistrationSagaConfig

// AuthServiceConfig aliases the outbound auth gRPC client configuration.
type AuthServiceConfig = config.AuthServiceConfig
