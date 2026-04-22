package gateway

import "errors"

var (
	ErrInvalidEmail              = errors.New("invalid email")
	ErrInvalidPassword           = errors.New("invalid password")
	ErrInvalidUsername           = errors.New("invalid username")
	ErrFailedToStartRegistration = errors.New("failed to start registration")
)

// RegistrationConflictError reports that registration cannot continue because
// either an in-progress session already exists or the owned username/email is
// already bound by the underlying services.
type RegistrationConflictError struct {
	State         string
	UsernameTaken bool
	EmailTaken    bool
}

// Error implements the error interface for registration conflict responses.
func (e *RegistrationConflictError) Error() string {
	return "registration conflict"
}
