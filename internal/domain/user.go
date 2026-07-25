package gateway

// User represents the public detailed freelancer profile returned by the
// gateway.
type User struct {
	UserID         string         `json:"user_id"`
	Username       string         `json:"username"`
	DisplayName    string         `json:"display_name"`
	AvatarID       string         `json:"avatar_id"`
	AvatarURL      string         `json:"avatar_url"`
	About          string         `json:"about"`
	ReviewsSummary *ReviewSummary `json:"reviews_summary"`
}

// GetUserProfileRequest loads a public user profile page by username with
// independent cursors for the gig and review sections.
type GetUserProfileRequest struct {
	Username      string `json:"username"`
	GigsCursor    string `json:"gigs_cursor"`
	ReviewsCursor string `json:"reviews_cursor"`
}

// UserProfile is the composite public payload returned by the profile page.
type UserProfile struct {
	User    *User           `json:"user"`
	Gigs    *GigPreviewList `json:"gigs"`
	Reviews *ReviewList     `json:"reviews"`
}
