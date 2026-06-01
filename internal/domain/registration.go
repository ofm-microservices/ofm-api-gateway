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

// AuthTokensResult is the public access/refresh token payload returned by the
// auth endpoints.
type AuthTokensResult struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// CompleteRegistrationResult is retained as a compatibility alias for the
// registration flow.
type CompleteRegistrationResult = AuthTokensResult

// SignInRequest is the public gateway payload used to sign into an existing
// auth credential.
type SignInRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

// RefreshTokensRequest is the public gateway payload used to rotate a refresh
// token.
type RefreshTokensRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// SignOutRequest is the public gateway payload used to revoke a refresh token
// and end the current session.
type SignOutRequest = RefreshTokensRequest

// RegistrationStatus is the gateway view of saga completion state.
type RegistrationStatus struct {
	SessionID string
	ClientID  string
	UserID    string
	Status    string
}
