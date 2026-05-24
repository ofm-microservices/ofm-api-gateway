package gateway

// CreateOrderRequest starts the public order creation flow.
type CreateOrderRequest struct {
	BuyerID              string `json:"buyer_id,omitempty"`
	BuyerEmail           string `json:"buyer_email,omitempty"`
	GigID                string `json:"gig_id"`
	PackageID            string `json:"package_id"`
	RealtimeConnectionID string `json:"realtime_connection_id,omitempty"`
	IdempotencyKey       string `json:"idempotency_key,omitempty"`
	RequestedAt          string `json:"requested_at,omitempty"`
}

// CreateOrderResult reports that the gateway accepted the order command.
type CreateOrderResult struct {
	SagaID      string          `json:"saga_id"`
	OrderID     string          `json:"order_id"`
	Status      string          `json:"status"`
	CheckoutURL string          `json:"checkout_url,omitempty"`
	Snapshot    *OrderSnapshot  `json:"snapshot,omitempty"`
	Questions   []OrderQuestion `json:"questions,omitempty"`
}

// ConfirmOrderRequest finalizes the order checkout flow and creates a payment
// session.
type ConfirmOrderRequest struct {
	OrderID              string `json:"order_id"`
	BuyerID              string `json:"buyer_id,omitempty"`
	RealtimeConnectionID string `json:"realtime_connection_id,omitempty"`
	IdempotencyKey       string `json:"idempotency_key,omitempty"`
	RequestedAt          string `json:"requested_at,omitempty"`
}

// ConfirmOrderResult reports that checkout can proceed after confirmation.
type ConfirmOrderResult struct {
	SagaID      string `json:"saga_id"`
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	CheckoutURL string `json:"checkout_url,omitempty"`
	PaymentID   string `json:"payment_id,omitempty"`
}

// DeliverOrderRequest submits the seller delivery payload.
type DeliverOrderRequest struct {
	OrderID         string   `json:"order_id"`
	SellerID        string   `json:"seller_id,omitempty"`
	DeliveryMessage string   `json:"delivery_message,omitempty"`
	AttachmentIDs   []string `json:"attachment_ids,omitempty"`
	RequestedAt     string   `json:"requested_at,omitempty"`
}

// DeliverOrderResult reports the delivery command outcome.
type DeliverOrderResult struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	CurrentStep string `json:"current_step,omitempty"`
}

// AcceptDeliveryRequest confirms the buyer accepted the seller delivery.
type AcceptDeliveryRequest struct {
	OrderID     string `json:"order_id"`
	BuyerID     string `json:"buyer_id,omitempty"`
	RequestedAt string `json:"requested_at,omitempty"`
}

// AcceptDeliveryResult reports the acceptance outcome.
type AcceptDeliveryResult struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	CurrentStep string `json:"current_step,omitempty"`
}

// RequestRevisionRequest asks the seller for revisions.
type RequestRevisionRequest struct {
	OrderID     string `json:"order_id"`
	BuyerID     string `json:"buyer_id,omitempty"`
	Reason      string `json:"reason,omitempty"`
	RequestedAt string `json:"requested_at,omitempty"`
}

// RequestRevisionResult reports the revision request outcome.
type RequestRevisionResult struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	CurrentStep string `json:"current_step,omitempty"`
}

// OpenDisputeRequest opens a buyer dispute for the current delivery.
type OpenDisputeRequest struct {
	OrderID     string `json:"order_id"`
	BuyerID     string `json:"buyer_id,omitempty"`
	Reason      string `json:"reason,omitempty"`
	RequestedAt string `json:"requested_at,omitempty"`
}

// OpenDisputeResult reports the dispute outcome.
type OpenDisputeResult struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	CurrentStep string `json:"current_step,omitempty"`
}

// OrderSnapshot captures the immutable commercial order data shown to the
// buyer before checkout.
type OrderSnapshot struct {
	GigID              string          `json:"gig_id"`
	PackageID          string          `json:"package_id"`
	SellerUserID       string          `json:"seller_user_id"`
	GigTitle           string          `json:"gig_title"`
	PackageTitle       string          `json:"package_title"`
	PackageDescription string          `json:"package_description"`
	PriceCents         int64           `json:"price_cents"`
	Currency           string          `json:"currency"`
	DeliveryDays       int32           `json:"delivery_days"`
	RevisionCount      int32           `json:"revision_count"`
	Questions          []OrderQuestion `json:"questions,omitempty"`
}

// OrderQuestion describes one snapshot question in the order start response.
type OrderQuestion struct {
	ID        string                `json:"id"`
	Text      string                `json:"text"`
	Type      string                `json:"type"`
	Required  bool                  `json:"required"`
	Options   []OrderQuestionOption `json:"options,omitempty"`
	SortOrder int32                 `json:"sort_order"`
}

// OrderQuestionOption describes one selectable option for a question.
type OrderQuestionOption struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Value string `json:"value"`
}

// SubmitOrderRequirementsRequest carries the buyer requirement answers.
type SubmitOrderRequirementsRequest struct {
	OrderID     string                   `json:"order_id"`
	BuyerID     string                   `json:"buyer_id,omitempty"`
	Answers     []OrderRequirementAnswer `json:"answers,omitempty"`
	RequestedAt string                   `json:"requested_at,omitempty"`
}

// OrderRequirementAnswer represents one question answer.
type OrderRequirementAnswer struct {
	QuestionID string `json:"question_id"`
	Value      string `json:"value"`
}

// SubmitOrderRequirementsResult reports the step outcome.
type SubmitOrderRequirementsResult struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	CurrentStep string `json:"current_step,omitempty"`
}

// SubmitOrderMessageRequest carries the buyer initial message.
type SubmitOrderMessageRequest struct {
	OrderID     string `json:"order_id"`
	BuyerID     string `json:"buyer_id,omitempty"`
	Message     string `json:"message"`
	RequestedAt string `json:"requested_at,omitempty"`
}

// SubmitOrderMessageResult reports the message step outcome.
type SubmitOrderMessageResult struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	CurrentStep string `json:"current_step,omitempty"`
}

// CreateOrderAttachmentUploadURLRequest asks for a presigned upload URL.
type CreateOrderAttachmentUploadURLRequest struct {
	OrderID     string `json:"order_id"`
	BuyerID     string `json:"buyer_id,omitempty"`
	FileName    string `json:"file_name"`
	MimeType    string `json:"mime_type"`
	SizeBytes   int64  `json:"size_bytes"`
	RequestedAt string `json:"requested_at,omitempty"`
}

// CreateOrderAttachmentUploadURLResult returns the upload URL.
type CreateOrderAttachmentUploadURLResult struct {
	OrderID      string `json:"order_id"`
	AttachmentID string `json:"attachment_id"`
	UploadURL    string `json:"upload_url"`
	FileKey      string `json:"file_key"`
}

// CompleteOrderAttachmentUploadRequest completes the attachment metadata.
type CompleteOrderAttachmentUploadRequest struct {
	OrderID      string `json:"order_id"`
	BuyerID      string `json:"buyer_id,omitempty"`
	AttachmentID string `json:"attachment_id"`
	FileKey      string `json:"file_key"`
	RequestedAt  string `json:"requested_at,omitempty"`
}

// CompleteOrderAttachmentUploadResult reports attachment completion.
type CompleteOrderAttachmentUploadResult struct {
	OrderID      string `json:"order_id"`
	AttachmentID string `json:"attachment_id"`
	Status       string `json:"status"`
}

// StartFreelancerOnboardingRequest starts the Stripe Connect onboarding flow.
type StartFreelancerOnboardingRequest struct {
	UserID     string `json:"user_id,omitempty"`
	Email      string `json:"email,omitempty"`
	Country    string `json:"country,omitempty"`
	ReturnURL  string `json:"return_url,omitempty"`
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
