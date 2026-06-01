package gateway

const (
	TierBasic    = "basic"
	TierStandard = "standard"
	TierPremium  = "premium"
)

// Gig represents the public draft and publish payload for the gig workflow.
type Gig struct {
	GigID                 string         `json:"gig_id,omitempty"`
	FreelancerID          string         `json:"-"`
	Freelancer            *User          `json:"freelancer,omitempty"`
	Slug                  string         `json:"slug,omitempty"`
	Title                 string         `json:"title,omitempty"`
	Description           string         `json:"description,omitempty"`
	CategoryID            int64          `json:"category_id,omitempty"`
	Currency              string         `json:"currency,omitempty"`
	Status                string         `json:"status,omitempty"`
	BasicInfoCompleted    bool           `json:"basic_info_completed,omitempty"`
	PackagesCompleted     bool           `json:"packages_completed,omitempty"`
	RequirementsCompleted bool           `json:"requirements_completed,omitempty"`
	MediaCompleted        bool           `json:"media_completed,omitempty"`
	PictureFileID         string         `json:"picture_file_id,omitempty"`
	PictureURL            string         `json:"picture_url,omitempty"`
	PublishedAt           string         `json:"published_at,omitempty"`
	CreatedAt             string         `json:"created_at,omitempty"`
	UpdatedAt             string         `json:"updated_at,omitempty"`
	Packages              []GigPackage   `json:"packages,omitempty"`
	Questions             []GigQuestion  `json:"questions,omitempty"`
	Media                 []GigMedia     `json:"media,omitempty"`
	Reviews               *ReviewList    `json:"reviews,omitempty"`
	ReviewsSummary        *ReviewSummary `json:"reviews_summary"`
}

// GigPackage is the public representation of one gig pricing tier.
type GigPackage struct {
	ID           string `json:"id,omitempty"`
	GigID        string `json:"gig_id,omitempty"`
	Tier         string `json:"tier"`
	Description  string `json:"description"`
	DeliveryDays int32  `json:"delivery_days"`
	PriceCents   int64  `json:"price_cents"`
	SortOrder    int32  `json:"sort_order,omitempty"`
}

// GigQuestion describes one buyer requirement question.
type GigQuestion struct {
	ID        string `json:"id,omitempty"`
	GigID     string `json:"gig_id,omitempty"`
	Content   string `json:"content"`
	SortOrder int32  `json:"sort_order,omitempty"`
}

// GigMedia describes one gig media file reference returned by gig-service.
type GigMedia struct {
	ID        string `json:"id,omitempty"`
	GigID     string `json:"gig_id,omitempty"`
	FileID    string `json:"file_id,omitempty"`
	URL       string `json:"url,omitempty"`
	SortOrder int32  `json:"sort_order,omitempty"`
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

// PublishGigRequest publishes a complete gig draft.
type PublishGigRequest struct {
	GigID        string `json:"gig_id"`
	FreelancerID string `json:"freelancer_id"`
}
