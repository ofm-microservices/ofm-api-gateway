package gateway

const (
	TierBasic    = "basic"
	TierStandard = "standard"
	TierPremium  = "premium"
)

// Gig represents the public draft and publish payload for the gig workflow.
type Gig struct {
	GigID                 string         `json:"gig_id"`
	FreelancerID          string         `json:"-"`
	Freelancer            *User          `json:"freelancer"`
	SellerUsername        string         `json:"seller_username"`
	Slug                  string         `json:"slug"`
	Title                 string         `json:"title"`
	Description           string         `json:"description"`
	ShortInfo             string         `json:"short_info"`
	CategoryID            int64          `json:"category_id"`
	Currency              string         `json:"currency"`
	Status                string         `json:"status"`
	BasicInfoCompleted    bool           `json:"basic_info_completed"`
	PackagesCompleted     bool           `json:"packages_completed"`
	RequirementsCompleted bool           `json:"requirements_completed"`
	MediaCompleted        bool           `json:"media_completed"`
	PictureFileID         string         `json:"picture_file_id"`
	PictureURL            string         `json:"picture_url"`
	PublishedAt           string         `json:"published_at"`
	CreatedAt             string         `json:"created_at"`
	UpdatedAt             string         `json:"updated_at"`
	Packages              []GigPackage   `json:"packages"`
	Questions             []GigQuestion  `json:"questions"`
	Media                 []GigMedia     `json:"media"`
	Reviews               *ReviewList    `json:"reviews"`
	ReviewsSummary        *ReviewSummary `json:"reviews_summary"`
}

// GigPackage is the public representation of one gig pricing tier.
type GigPackage struct {
	ID           string `json:"id"`
	GigID        string `json:"gig_id"`
	Tier         string `json:"tier"`
	Description  string `json:"description"`
	DeliveryDays int32  `json:"delivery_days"`
	PriceCents   int64  `json:"price_cents"`
	SortOrder    int32  `json:"sort_order"`
}

// GigQuestion describes one buyer requirement question.
type GigQuestion struct {
	ID        string `json:"id"`
	GigID     string `json:"gig_id"`
	Content   string `json:"content"`
	SortOrder int32  `json:"sort_order"`
}

// GigMedia describes one gig media file reference returned by gig-service.
type GigMedia struct {
	ID        string `json:"id"`
	GigID     string `json:"gig_id"`
	FileID    string `json:"file_id"`
	URL       string `json:"url"`
	SortOrder int32  `json:"sort_order"`
}

// GigMediaUpload carries one uploaded media file through the gateway.
type GigMediaUpload struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Data        []byte `json:"data"`
}

// CreateGigDraftRequest starts a new server-side gig draft.
type CreateGigDraftRequest struct {
	FreelancerID string `json:"freelancer_id"`
}

// UpdateGigBasicInfoRequest updates the public metadata of a gig draft.
type UpdateGigBasicInfoRequest struct {
	GigID        string `json:"gig_id"`
	FreelancerID string `json:"freelancer_id"`
	Title        string `json:"title"`
	ShortInfo    string `json:"short_info"`
	Description  string `json:"description"`
	CategoryID   int64  `json:"category_id"`
	Currency     string `json:"currency"`
}

// ReplaceGigPackagesRequest replaces all draft packages.
type ReplaceGigPackagesRequest struct {
	GigID        string       `json:"gig_id"`
	FreelancerID string       `json:"freelancer_id"`
	Packages     []GigPackage `json:"packages"`
}

// ReplaceGigQuestionsRequest replaces all draft questions.
type ReplaceGigQuestionsRequest struct {
	GigID        string        `json:"gig_id"`
	FreelancerID string        `json:"freelancer_id"`
	Questions    []GigQuestion `json:"questions"`
}

// ReplaceGigMediaRequest uploads the ordered media batch for a gig draft.
type ReplaceGigMediaRequest struct {
	GigID        string           `json:"gig_id"`
	FreelancerID string           `json:"freelancer_id"`
	Files        []GigMediaUpload `json:"files"`
}

// GetGigDraftRequest loads the current draft state.
type GetGigDraftRequest struct {
	GigID        string `json:"gig_id"`
	FreelancerID string `json:"freelancer_id"`
}

// GetGigBySlugRequest loads the public gig view by slug.
type GetGigBySlugRequest struct {
	Username string `json:"username"`
	Slug     string `json:"slug"`
	Cursor   string `json:"cursor"`
}

// GetPreviewGigsByFreelancerUsernameRequest loads the public freelancer gig
// preview list by username.
type GetPreviewGigsByFreelancerUsernameRequest struct {
	Username string `json:"username"`
	Page     int32  `json:"page"`
	Limit    int32  `json:"limit"`
}

// GigPreview represents one item in the freelancer preview list.
type GigPreview struct {
	GigID             string `json:"gig_id"`
	Slug              string `json:"slug"`
	Title             string `json:"title"`
	ShortInfo         string `json:"short_info"`
	MinimumPriceCents int64  `json:"minimum_price_cents"`
	PictureURL        string `json:"picture_url"`
	CreatedAt         string `json:"created_at"`
	Status            string `json:"status"`
	PublishedAt       string `json:"published_at"`
	UpdatedAt         string `json:"updated_at"`
	RatingAvg         float64 `json:"rating_avg"`
	TotalReviews      int64   `json:"total_reviews"`
	OrderCount        int64   `json:"order_count"`
}

// GigPreviewList wraps one page of freelancer preview gigs.
type GigPreviewList struct {
	Items      []GigPreview `json:"items"`
	Page       int32        `json:"page"`
	Limit      int32        `json:"limit"`
	TotalPages int32        `json:"total_pages"`
}

// GetMyGigsRequest loads the authenticated owner gig list.
type GetMyGigsRequest struct {
	UserID string `json:"user_id"`
	Status string `json:"status"`
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Page   int32  `json:"page"`
	Limit  int32  `json:"limit"`
}

// GigPreviewPage wraps one page of owner gig previews.
type GigPreviewPage struct {
	Items      []GigPreview `json:"items"`
	Page       int32        `json:"page"`
	Limit      int32        `json:"limit"`
	TotalPages int32        `json:"total_pages"`
}

// PublishGigRequest publishes a complete gig draft.
type PublishGigRequest struct {
	GigID        string `json:"gig_id"`
	FreelancerID string `json:"freelancer_id"`
	Username     string `json:"username"`
}
