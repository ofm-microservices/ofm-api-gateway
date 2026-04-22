package service

import "errors"

var (
	ErrNilRegistrationClient = errors.New("registration client is nil")
	ErrNilLogger             = errors.New("logger is nil")
)
