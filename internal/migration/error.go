package migration

import "errors"

var (
	ErrEmptySessionID = errors.New("registration session id is empty")
	ErrNilRedis       = errors.New("redis client is nil")
)
