package gateway

// SearchRequest carries the public search query parameters.
type SearchRequest struct {
	Query  string
	Sort   int32
	Order  int32
	Cursor string
}

// SearchResult describes one public gig search result.
type SearchResult struct {
	ID             string  `json:"id"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	Picture        string  `json:"picture"`
	ReviewsCount   int64   `json:"reviews_count"`
	Rating         float64 `json:"rating"`
	MinPrice       int64   `json:"min_price"`
	Slug           string  `json:"slug"`
	FreelancerID   string  `json:"freelancer_id"`
	SellerUsername string  `json:"seller_username"`
	PublishedAt    string  `json:"published_at"`
}

// SearchResponse reports one page of public search results.
type SearchResponse struct {
	Items   []SearchResult `json:"items"`
	Cursor  string         `json:"cursor"`
	HasMore bool           `json:"has_more"`
}
