package http

// APIErrorResponse describes the shared error payload returned by the HTTP
// gateway endpoints.
type APIErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}
