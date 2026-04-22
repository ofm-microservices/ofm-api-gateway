package grpc

import "errors"

var (
	ErrEmptyAddress = errors.New("registration saga address is empty")
	ErrNilLogger    = errors.New("logger is nil")
)
