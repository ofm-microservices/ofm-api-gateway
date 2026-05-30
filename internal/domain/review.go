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
