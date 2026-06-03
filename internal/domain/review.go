package gateway

// CreateReviewRequest submits a buyer-authenticated review for a completed order.
type CreateReviewRequest struct {
	OrderID       string `json:"order_id"`
	Content       string `json:"content"`
	Rating        int32  `json:"rating"`
	BuyerID       string `json:"-"`
	BuyerUsername string `json:"-"`
	RequestedAt   string `json:"requested_at"`
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
	Cursor string `json:"cursor"`
}

// GetGigReviewsResult contains the public review list for one gig.
type GetGigReviewsResult struct {
	Reviews *ReviewList `json:"reviews"`
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

// GetReviewsBySellerUsernameRequest loads public reviews for one seller by username.
type GetReviewsBySellerUsernameRequest struct {
	Username string `json:"username"`
	Cursor   string `json:"cursor"`
}

// GetReviewsBySellerUsernameResult contains the public review list for one seller.
type GetReviewsBySellerUsernameResult struct {
	Reviews *ReviewList `json:"reviews"`
}

// ListSellerReviewsRequest is a compatibility alias for the seller-username request.
type ListSellerReviewsRequest = GetReviewsBySellerUsernameRequest

// ListSellerReviewsResult is a compatibility alias for the seller-username result.
type ListSellerReviewsResult = GetReviewsBySellerUsernameResult

// GetUserRatingSummaryByUsernameRequest loads the user rating summary for one
// public username handle.
type GetUserRatingSummaryByUsernameRequest struct {
	Username string `json:"username"`
}

// ReviewAuthor describes the user metadata attached to a public review.
type ReviewAuthor struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

// Review represents one public gig review returned by api-gateway.
type Review struct {
	ReviewID       string        `json:"review_id"`
	OrderID        string        `json:"order_id"`
	GigID          string        `json:"gig_id"`
	BuyerUserID    string        `json:"buyer_user_id"`
	Content        string        `json:"content"`
	Rating         int32         `json:"rating"`
	CreatedAt      string        `json:"created_at"`
	SellerUsername string        `json:"seller_username"`
	Author         *ReviewAuthor `json:"author"`
}

// ReviewSummary aggregates public rating counts for a gig.
type ReviewSummary struct {
	RatingAvg    float64 `json:"rating_avg"`
	TotalReviews int64   `json:"total_reviews"`
	Stars5       int64   `json:"stars_5"`
	Stars4       int64   `json:"stars_4"`
	Stars3       int64   `json:"stars_3"`
	Stars2       int64   `json:"stars_2"`
	Stars1       int64   `json:"stars_1"`
}
