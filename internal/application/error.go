package service

import "errors"

var (
	ErrNilRegistrationClient      = errors.New("registration client is nil")
	ErrNilTokenIssuer             = errors.New("token issuer is nil")
	ErrNilGigClient               = errors.New("gig client is nil")
	ErrNilOrderClient             = errors.New("order client is nil")
	ErrNilPaymentOnboardingClient = errors.New("payment onboarding client is nil")
	ErrNilLogger                  = errors.New("logger is nil")
)
