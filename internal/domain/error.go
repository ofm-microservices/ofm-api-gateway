package gateway

import "errors"

var (
	ErrInvalidEmail                 = errors.New("invalid email")
	ErrInvalidPassword              = errors.New("invalid password")
	ErrInvalidUsername              = errors.New("invalid username")
	ErrInvalidSessionID             = errors.New("invalid session id")
	ErrInvalidClientID              = errors.New("invalid client id")
	ErrInvalidVerificationCode      = errors.New("invalid verification code")
	ErrRegistrationNotCompleted     = errors.New("registration is not completed")
	ErrRegistrationAlreadyClaimed   = errors.New("registration tokens already claimed")
	ErrFailedToStartRegistration    = errors.New("failed to start registration")
	ErrFailedToVerifyEmail          = errors.New("failed to verify email")
	ErrFailedToCompleteRegistration = errors.New("failed to complete registration")
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
