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
}

// RegistrationService validates the public signup request and delegates
// orchestration to registration-saga-service.
type RegistrationService interface {
	// SignUp validates the public request and starts the registration saga.
	SignUp(ctx context.Context, req gateway.SignUpRequest) (*gateway.SignUpResult, error)
}

// Logger aliases the shared structured logger used by the application layer.
type Logger = logging.Logger
