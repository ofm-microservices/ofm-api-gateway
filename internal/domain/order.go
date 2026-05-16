package gateway

// CreateOrderRequest starts the public order creation flow.
type CreateOrderRequest struct {
	SagaID               string `json:"saga_id,omitempty"`
	OrderID              string `json:"order_id,omitempty"`
	RequestedAt          string `json:"requested_at,omitempty"`
	IdempotencyKey       string `json:"idempotency_key,omitempty"`
	BuyerID              string `json:"buyer_id,omitempty"`
	BuyerEmail           string `json:"buyer_email"`
	RealtimeConnectionID string `json:"realtime_connection_id,omitempty"`
	GigID                string `json:"gig_id"`
	GigTitle             string `json:"gig_title"`
	PackageID            string `json:"package_id"`
	PackageTier          string `json:"package_tier"`
	PackageDescription   string `json:"package_description"`
	PackageDeliveryDays  int32  `json:"package_delivery_days"`
	PriceCents           int64  `json:"price_cents"`
	Currency             string `json:"currency"`
}

// CreateOrderResult reports that the gateway accepted the order command.
type CreateOrderResult struct {
	SagaID  string `json:"saga_id"`
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

// StartFreelancerOnboardingRequest starts the Stripe Connect onboarding flow.
type StartFreelancerOnboardingRequest struct {
	UserID    string `json:"user_id,omitempty"`
	Email     string `json:"email,omitempty"`
	Country   string `json:"country,omitempty"`
	ReturnURL string `json:"return_url,omitempty"`
	RefreshURL string `json:"refresh_url,omitempty"`
}

// StartFreelancerOnboardingResult reports the Stripe onboarding URL and
// connected-account state.
type StartFreelancerOnboardingResult struct {
	UserID           string `json:"user_id"`
	StripeAccountID  string `json:"stripe_account_id"`
	OnboardingURL    string `json:"onboarding_url,omitempty"`
	Status           string `json:"status"`
	DetailsSubmitted bool   `json:"details_submitted"`
	ChargesEnabled   bool   `json:"charges_enabled"`
	PayoutsEnabled   bool   `json:"payouts_enabled"`
	DisabledReason   string `json:"disabled_reason,omitempty"`
}
