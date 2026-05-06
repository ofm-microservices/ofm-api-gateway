package service

import (
	"api-gateway/internal/domain"
	"context"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"
)

// RegistrationPublisher is the outbound boundary used to start registration in
// another service.
type RegistrationPublisher interface {
	// StartRegistration delegates registration start to the saga service.
	StartRegistration(ctx context.Context, req gateway.SignUpRequest) (*gateway.SignUpResult, error)
	// VerifyEmail delegates email verification to the saga service.
	VerifyEmail(ctx context.Context, req gateway.VerifyEmailRequest) (*gateway.VerifyEmailResult, error)
	// GetRegistrationStatus returns saga state before token exchange.
	GetRegistrationStatus(ctx context.Context, sessionID, clientID string) (*gateway.RegistrationStatus, error)
}

// TokenIssuer is the outbound auth-service boundary for final token exchange.
type TokenIssuer interface {
	// IssueRegistrationTokens creates login tokens after saga completion.
	IssueRegistrationTokens(ctx context.Context, userID string) (*gateway.CompleteRegistrationResult, error)
	Close() error
}

// RegistrationService validates the public signup request and delegates
// orchestration to registration-saga-service.
type RegistrationService interface {
	// SignUp validates the public request and starts the registration saga.
	SignUp(ctx context.Context, req gateway.SignUpRequest) (*gateway.SignUpResult, error)
	// VerifyEmail accepts the emailed code and starts saga completion.
	VerifyEmail(ctx context.Context, req gateway.VerifyEmailRequest) (*gateway.VerifyEmailResult, error)
	// CompleteRegistration exchanges a completed saga for auth-owned tokens.
	CompleteRegistration(ctx context.Context, req gateway.CompleteRegistrationRequest) (*gateway.CompleteRegistrationResult, error)
}

// Logger aliases the shared structured logger used by the application layer.
type Logger = logging.Logger
