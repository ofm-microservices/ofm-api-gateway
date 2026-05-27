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
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Picture      string  `json:"picture"`
	ReviewsCount int64   `json:"reviewsCount"`
	Rating       float64 `json:"rating"`
	MinPrice     int64   `json:"minPrice"`
	Slug         string  `json:"slug"`
	FreelancerID string  `json:"freelancer_id"`
	PublishedAt  string  `json:"published_at"`
}

// SearchResponse reports one page of public search results.
type SearchResponse struct {
	Services []SearchResult `json:"services"`
	Cursor   string         `json:"cursor,omitempty"`
	HasMore  bool           `json:"hasMore"`
}
