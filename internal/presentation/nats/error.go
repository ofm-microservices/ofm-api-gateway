package nats

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyNATSURL = errors.New("nats url is empty")
	ErrEmptySubject = errors.New("subject is empty")
	ErrNilLogger    = errors.New("logger is nil")
)

// WrapConnectToNATSError annotates NATS connection failures.
func WrapConnectToNATSError(err error) error {
	return fmt.Errorf("connect to nats: %w", err)
}

// WrapMarshalEventError annotates signup payload serialization failures.
func WrapMarshalEventError(err error) error {
	return fmt.Errorf("marshal event: %w", err)
}

// WrapPublishToNATSError annotates publish failures for the given subject.
func WrapPublishToNATSError(subject string, err error) error {
	return fmt.Errorf("publish to nats (%s): %w", subject, err)
}

// WrapFlushNATSError annotates publisher flush failures.
func WrapFlushNATSError(err error) error {
	return fmt.Errorf("flush nats publisher: %w", err)
}
