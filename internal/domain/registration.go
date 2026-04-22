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
