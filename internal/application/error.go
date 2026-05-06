package service

import "errors"

var (
	ErrNilRegistrationClient = errors.New("registration client is nil")
	ErrNilTokenIssuer        = errors.New("token issuer is nil")
	ErrNilLogger             = errors.New("logger is nil")
)
