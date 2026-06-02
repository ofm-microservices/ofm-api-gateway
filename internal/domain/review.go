package gateway

// CreateReviewRequest submits a buyer-authenticated review for a completed order.
type CreateReviewRequest struct {
	OrderID     string `json:"order_id"`
	Content     string `json:"content"`
	Rating      int32  `json:"rating"`
	BuyerID     string `json:"-"`
	RequestedAt string `json:"requested_at,omitempty"`
}

// CreateReviewResult reports the created review.
type CreateReviewResult struct {
	ReviewID  string `json:"review_id"`
	OrderID   string `json:"order_id"`
	GigID     string `json:"gig_id"`
	BuyerID   string `json:"buyer_id"`
	Content   string `json:"content"`
	Rating    int32  `json:"rating"`
	CreatedAt string `json:"created_at"`
}

// GetGigReviewsRequest loads public reviews for one gig.
type GetGigReviewsRequest struct {
	GigID  string `json:"gig_id"`
	Cursor string `json:"cursor,omitempty"`
}

// GetGigReviewsResult contains the public review list for one gig.
type GetGigReviewsResult struct {
	Reviews *ReviewList `json:"reviews,omitempty"`
}

// ReviewList wraps one page of public gig reviews.
type ReviewList struct {
	Items   []Review `json:"items"`
	Cursor  string   `json:"cursor"`
	HasMore bool     `json:"has_more"`
}

// GetGigReviewsSummaryRequest loads the aggregate rating summary for one gig.
type GetGigReviewsSummaryRequest struct {
	GigID string `json:"gig_id"`
}

// GetUserRatingSummaryByUsernameRequest loads the user rating summary for one
// public username handle.
type GetUserRatingSummaryByUsernameRequest struct {
	Username string `json:"username"`
}

// ReviewAuthor describes the user metadata attached to a public review.
type ReviewAuthor struct {
	UserID      string `json:"user_id,omitempty"`
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	AvatarURL   string `json:"avatar_url"`
}

// Review represents one public gig review returned by api-gateway.
type Review struct {
	ReviewID    string        `json:"review_id,omitempty"`
	OrderID     string        `json:"order_id,omitempty"`
	GigID       string        `json:"gig_id,omitempty"`
	BuyerUserID string        `json:"buyer_user_id,omitempty"`
	Content     string        `json:"content,omitempty"`
	Rating      int32         `json:"rating,omitempty"`
	CreatedAt   string        `json:"created_at,omitempty"`
	Author      *ReviewAuthor `json:"author,omitempty"`
}

// ReviewSummary aggregates public rating counts for a gig.
type ReviewSummary struct {
	RatingAvg    float64 `json:"rating_avg,omitempty"`
	TotalReviews int64   `json:"total_reviews,omitempty"`
	Stars5       int64   `json:"stars_5,omitempty"`
	Stars4       int64   `json:"stars_4,omitempty"`
	Stars3       int64   `json:"stars_3,omitempty"`
	Stars2       int64   `json:"stars_2,omitempty"`
	Stars1       int64   `json:"stars_1,omitempty"`
}
