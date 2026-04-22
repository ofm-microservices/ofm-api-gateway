package grpc

import (
	"api-gateway/config"
	gateway "api-gateway/internal/domain"
	"context"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"
	registrationv1 "github.com/ofm-microseervices/ofm-common/proto/registration/v1"
)

// Logger aliases the shared logger contract used by the gRPC adapter.
type Logger = logging.Logger

// SignUpRequest aliases the gateway-domain signup request transported over gRPC.
type SignUpRequest = gateway.SignUpRequest

// SignUpResult aliases the gateway-domain signup result returned from the saga.
type SignUpResult = gateway.SignUpResult

// Client is the gateway-facing gRPC adapter for registration-saga-service.
type Client interface {
	StartRegistration(ctx context.Context, req SignUpRequest) (*SignUpResult, error)
	Close() error
}

// RegistrationMapper translates between gateway-domain signup types and the
// shared registration gRPC contract.
type RegistrationMapper interface {
	ToStartRegistrationRequest(req SignUpRequest) *registrationv1.StartRegistrationRequest
	ToSignUpResult(res *registrationv1.StartRegistrationResponse) *SignUpResult
	ToStartRegistrationError(err error) error
}

// RegistrationSagaConfig aliases the outbound saga gRPC client configuration.
type RegistrationSagaConfig = config.RegistrationSagaConfig
