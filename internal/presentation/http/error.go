package http

import "errors"

var (
	ErrNilRegistrationService      = errors.New("registration service is nil")
	ErrNilAuthSessionService       = errors.New("auth session service is nil")
	ErrNilAuthMeService            = errors.New("auth me service is nil")
	ErrNilGigService               = errors.New("gig service is nil")
	ErrNilOrderService             = errors.New("order service is nil")
	ErrNilReviewService            = errors.New("review service is nil")
	ErrNilSearchService            = errors.New("search service is nil")
	ErrNilPaymentOnboardingService = errors.New("payment onboarding service is nil")
	ErrNilLogger                   = errors.New("logger is nil")
	ErrNilAuthHandler              = errors.New("auth handler is nil")
	ErrNilGigHandler               = errors.New("gig handler is nil")
	ErrNilOrderHandler             = errors.New("order handler is nil")
	ErrNilReviewHandler            = errors.New("review handler is nil")
	ErrNilSearchHandler            = errors.New("search handler is nil")
	ErrInvalidRequestBody          = errors.New("invalid request body")
)
