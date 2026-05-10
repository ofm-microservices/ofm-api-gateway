package http

import "errors"

var (
	ErrNilRegistrationService = errors.New("registration service is nil")
	ErrNilGigService          = errors.New("gig service is nil")
	ErrNilLogger              = errors.New("logger is nil")
	ErrNilAuthHandler         = errors.New("auth handler is nil")
	ErrNilGigHandler          = errors.New("gig handler is nil")
	ErrInvalidRequestBody     = errors.New("invalid request body")
)
