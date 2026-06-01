package gateway

// User represents the public detailed freelancer profile returned by the
// gateway.
type User struct {
	UserID         string         `json:"user_id,omitempty"`
	Username       string         `json:"username,omitempty"`
	DisplayName    string         `json:"display_name,omitempty"`
	AvatarID       string         `json:"avatar_id,omitempty"`
	AvatarURL      string         `json:"avatar_url,omitempty"`
	About          string         `json:"about,omitempty"`
	ReviewsSummary *ReviewSummary `json:"reviews_summary,omitempty"`
}
