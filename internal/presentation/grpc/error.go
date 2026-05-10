package grpc

import "errors"

var (
	ErrEmptyAddress           = errors.New("registration saga address is empty")
	ErrEmptyGigServiceAddress = errors.New("gig service address is empty")
	ErrNilLogger              = errors.New("logger is nil")
)
