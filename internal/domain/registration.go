package gateway

// SignUpRequest is the public gateway payload used to start registration.
type SignUpRequest struct {
	ClientID  string `json:"client_id,omitempty"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	Surname   string `json:"surname"`
}

// SignUpResult is the immediate response returned after the saga session is
// created.
type SignUpResult struct {
	SessionID     string `json:"session_id,omitempty"`
	ClientID      string `json:"client_id,omitempty"`
	UserID        string `json:"user_id,omitempty"`
	Status        string `json:"status"`
	ConflictState string `json:"conflict_state,omitempty"`
	UsernameTaken bool   `json:"username_taken,omitempty"`
	EmailTaken    bool   `json:"email_taken,omitempty"`
}

// VerifyEmailRequest is the public gateway payload for submitting a
// registration email verification code.
type VerifyEmailRequest struct {
	SessionID string `json:"session_id"`
	ClientID  string `json:"client_id"`
	Code      string `json:"code"`
}

// VerifyEmailResult reports that the saga accepted email verification.
type VerifyEmailResult struct {
	SessionID string `json:"session_id"`
	ClientID  string `json:"client_id"`
	Status    string `json:"status"`
}

// CompleteRegistrationRequest asks api-gateway to exchange a completed saga for
// auth-owned login tokens.
type CompleteRegistrationRequest struct {
	SessionID string `json:"session_id"`
	ClientID  string `json:"client_id"`
}

// CompleteRegistrationResult is the public token response returned after saga
// completion.
type CompleteRegistrationResult struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// RegistrationStatus is the gateway view of saga completion state.
type RegistrationStatus struct {
	SessionID string
	ClientID  string
	UserID    string
	Status    string
}
