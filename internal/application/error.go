package service

import "errors"

var (
	ErrNilRegistrationClient      = errors.New("registration client is nil")
	ErrNilTokenIssuer             = errors.New("token issuer is nil")
	ErrNilAuthSessionClient       = errors.New("auth session client is nil")
	ErrNilGigClient               = errors.New("gig client is nil")
	ErrNilUserClient              = errors.New("user client is nil")
	ErrNilOrderClient             = errors.New("order client is nil")
	ErrNilReviewClient            = errors.New("review client is nil")
	ErrNilSearchClient            = errors.New("search client is nil")
	ErrNilPaymentOnboardingClient = errors.New("payment onboarding client is nil")
	ErrNilLogger                  = errors.New("logger is nil")
)
