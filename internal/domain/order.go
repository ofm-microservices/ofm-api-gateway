package gateway

import (
	"errors"
	"strings"
)

// CreateOrderRequest starts the public order creation flow.
type CreateOrderRequest struct {
	BuyerID        string `json:"buyer_id"`
	BuyerEmail     string `json:"buyer_email"`
	GigID          string `json:"gig_id"`
	PackageID      string `json:"package_id"`
	IdempotencyKey string `json:"idempotency_key"`
	RequestedAt    string `json:"requested_at"`
}

// CreateOrderResult reports that the gateway accepted the order command.
type CreateOrderResult struct {
	SagaID      string          `json:"saga_id"`
	OrderID     string          `json:"order_id"`
	Status      string          `json:"status"`
	CheckoutURL string          `json:"checkout_url"`
	Snapshot    *OrderSnapshot  `json:"snapshot"`
	Questions   []OrderQuestion `json:"questions"`
}

// ConfirmOrderRequest finalizes the order checkout flow and creates a payment
// session.
type ConfirmOrderRequest struct {
	OrderID        string `json:"order_id"`
	BuyerID        string `json:"buyer_id"`
	IdempotencyKey string `json:"idempotency_key"`
	RequestedAt    string `json:"requested_at"`
}

// ConfirmOrderResult reports that checkout can proceed after confirmation.
type ConfirmOrderResult struct {
	SagaID      string `json:"saga_id"`
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	CheckoutURL string `json:"checkout_url"`
	PaymentID   string `json:"payment_id"`
}

// DeliverOrderRequest submits the seller delivery payload.
type DeliverOrderRequest struct {
	OrderID         string   `json:"order_id"`
	SellerID        string   `json:"seller_id"`
	DeliveryMessage string   `json:"delivery_message"`
	AttachmentIDs   []string `json:"attachment_ids"`
	RequestedAt     string   `json:"requested_at"`
}

// DeliverOrderResult reports the delivery command outcome.
type DeliverOrderResult struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	CurrentStep string `json:"current_step"`
}

// AcceptDeliveryRequest confirms the buyer accepted the seller delivery.
type AcceptDeliveryRequest struct {
	OrderID     string `json:"order_id"`
	BuyerID     string `json:"buyer_id"`
	RequestedAt string `json:"requested_at"`
}

// AcceptDeliveryResult reports the acceptance outcome.
type AcceptDeliveryResult struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	CurrentStep string `json:"current_step"`
}

// RequestRevisionRequest asks the seller for revisions.
type RequestRevisionRequest struct {
	OrderID     string `json:"order_id"`
	BuyerID     string `json:"buyer_id"`
	Reason      string `json:"reason"`
	RequestedAt string `json:"requested_at"`
}

// RequestRevisionResult reports the revision request outcome.
type RequestRevisionResult struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	CurrentStep string `json:"current_step"`
}

// OpenDisputeRequest opens a dispute for the authenticated order owner.
type OpenDisputeRequest struct {
	OrderID     string `json:"order_id"`
	ActorID     string `json:"-"`
	Reason      string `json:"reason"`
	RequestedAt string `json:"requested_at"`
}

// OpenDisputeResult reports the dispute outcome.
type OpenDisputeResult struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	CurrentStep string `json:"current_step"`
}

// ResolveDisputeRequest lets an admin split the disputed settlement between
// the freelancer and the customer.
type ResolveDisputeRequest struct {
	OrderID              string `json:"order_id"`
	AdminUserID          string `json:"admin_user_id"`
	FreelancerPercentage int32  `json:"freelancer_percentage"`
	CustomerPercentage   int32  `json:"customer_percentage"`
	Reason               string `json:"reason"`
	RequestedAt          string `json:"requested_at"`
}

// ResolveDisputeResult reports the settlement outcome after admin resolution.
type ResolveDisputeResult struct {
	OrderID          string `json:"order_id"`
	Status           string `json:"status"`
	CurrentStep      string `json:"current_step"`
	PaymentReleaseID string `json:"payment_release_id"`
	StripeTransferID string `json:"stripe_transfer_id"`
	StripeRefundID   string `json:"stripe_refund_id"`
}

// ParticipantRole is the shared customer/freelancer view role used by user-scoped resources.
type ParticipantRole string

const (
	ParticipantRoleUnspecified ParticipantRole = ""
	ParticipantRoleCustomer    ParticipantRole = "customer"
	ParticipantRoleFreelancer  ParticipantRole = "freelancer"
	RoleAdmin                                  = "admin"
)

// ParseParticipantRole validates a shared view role from query parameters.
func ParseParticipantRole(raw string) (ParticipantRole, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case string(ParticipantRoleCustomer):
		return ParticipantRoleCustomer, nil
	case string(ParticipantRoleFreelancer):
		return ParticipantRoleFreelancer, nil
	default:
		return ParticipantRoleUnspecified, errors.New("invalid participant role")
	}
}

// GetOrderPreviewByIDRequest loads a user-scoped minimal order preview.
type GetOrderPreviewByIDRequest struct {
	OrderID string
	UserID  string
	Role    ParticipantRole
}

// GetOrderRequirementsByIDRequest loads the order requirements page for a
// user-scoped order.
type GetOrderRequirementsByIDRequest struct {
	OrderID string
	UserID  string
}

// GetOrderDeliveryByIDRequest loads the order delivery page for a
// user-scoped order.
type GetOrderDeliveryByIDRequest struct {
	OrderID string
	UserID  string
}

// OrderPreview is the minimal order payload returned by the preview endpoint.
type OrderPreview struct {
	OrderID   string `json:"order_id"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

// OrderPreviewGig describes the gig snapshot returned alongside the order preview.
type OrderPreviewGig struct {
	GigID               string `json:"gig_id"`
	Title               string `json:"title"`
	PictureURL          string `json:"picture_url"`
	PackageID           string `json:"package_id"`
	PackageTitle        string `json:"package_title"`
	PriceCents          int64  `json:"price_cents"`
	Currency            string `json:"currency"`
	Description         string `json:"description"`
	PackageDeliveryDays int32  `json:"package_delivery_days"`
}

// OrderPreviewUser describes one participant snapshot in the order preview.
type OrderPreviewUser struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

// OrderRequirementsQuestion describes one requirements question snapshot.
type OrderRequirementsQuestion struct {
	QuestionID string `json:"question_id"`
	Text       string `json:"text"`
	Type       string `json:"type"`
	Required   bool   `json:"required"`
	SortOrder  int32  `json:"sort_order"`
}

// OrderRequirementsAnswer describes one optional buyer answer snapshot.
type OrderRequirementsAnswer struct {
	Value string `json:"value"`
}

// OrderRequirementsQuestionAnswer describes a question and its optional answer.
type OrderRequirementsQuestionAnswer struct {
	Question *OrderRequirementsQuestion `json:"question"`
	Answer   *OrderRequirementsAnswer   `json:"answer"`
}

// OrderRequirementsCustomerMessage describes the buyer message snapshot.
type OrderRequirementsCustomerMessage struct {
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// OrderDelivery describes the seller delivery payload returned by the delivery endpoint.
type OrderDelivery struct {
	DeliveryMessage string `json:"delivery_message"`
	CreatedAt       string `json:"created_at"`
}

// OrderDeliveryFile describes one delivery attachment returned by the delivery endpoint.
type OrderDeliveryFile struct {
	FileID    string `json:"file_id"`
	FileURL   string `json:"file_url"`
	SortOrder int32  `json:"sort_order"`
	CreatedAt string `json:"created_at"`
}

// OrderPreviewPayment describes the payment snapshot returned alongside the order preview.
type OrderPreviewPayment struct {
	PaymentID   string `json:"payment_id"`
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// GetOrderPreviewByIDResult wraps the preview payload in the public HTTP response.
type GetOrderPreviewByIDResult struct {
	Order      *OrderPreview        `json:"order"`
	Gig        *OrderPreviewGig     `json:"gig"`
	Payment    *OrderPreviewPayment `json:"payment"`
	Customer   *OrderPreviewUser    `json:"customer"`
	Freelancer *OrderPreviewUser    `json:"freelancer"`
}

// GetOrderRequirementsByIDResult wraps the order requirements payload in the
// public HTTP response.
type GetOrderRequirementsByIDResult struct {
	QuestionsAnswers []OrderRequirementsQuestionAnswer `json:"questions_answers"`
	CustomerMessage  *OrderRequirementsCustomerMessage `json:"customer_message"`
}

// GetOrderDeliveryByIDResult wraps the seller delivery payload in the public HTTP response.
type GetOrderDeliveryByIDResult struct {
	OrderDelivery      *OrderDelivery      `json:"order_delivery"`
	OrderDeliveryFiles []OrderDeliveryFile `json:"order_delivery_files"`
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
	Questions          []OrderQuestion `json:"questions"`
}

// OrderQuestion describes one snapshot question in the order start response.
type OrderQuestion struct {
	ID        string                `json:"id"`
	Text      string                `json:"text"`
	Type      string                `json:"type"`
	Required  bool                  `json:"required"`
	Options   []OrderQuestionOption `json:"options"`
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
	BuyerID     string                   `json:"buyer_id"`
	Answers     []OrderRequirementAnswer `json:"answers"`
	RequestedAt string                   `json:"requested_at"`
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
	CurrentStep string `json:"current_step"`
}

// SubmitOrderMessageRequest carries the buyer initial message.
type SubmitOrderMessageRequest struct {
	OrderID     string `json:"order_id"`
	BuyerID     string `json:"buyer_id"`
	Message     string `json:"message"`
	RequestedAt string `json:"requested_at"`
}

// SubmitOrderMessageResult reports the message step outcome.
type SubmitOrderMessageResult struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	CurrentStep string `json:"current_step"`
}

// CreateOrderAttachmentUploadURLRequest asks for a presigned upload URL.
type CreateOrderAttachmentUploadURLRequest struct {
	OrderID     string `json:"order_id"`
	BuyerID     string `json:"buyer_id"`
	FileName    string `json:"file_name"`
	MimeType    string `json:"mime_type"`
	SizeBytes   int64  `json:"size_bytes"`
	RequestedAt string `json:"requested_at"`
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
	BuyerID      string `json:"buyer_id"`
	AttachmentID string `json:"attachment_id"`
	FileKey      string `json:"file_key"`
	RequestedAt  string `json:"requested_at"`
}

// CompleteOrderAttachmentUploadResult reports attachment completion.
type CompleteOrderAttachmentUploadResult struct {
	OrderID      string `json:"order_id"`
	AttachmentID string `json:"attachment_id"`
	Status       string `json:"status"`
}

// StartFreelancerOnboardingRequest starts the Stripe Connect onboarding flow.
type StartFreelancerOnboardingRequest struct {
	UserID     string `json:"user_id"`
	Email      string `json:"email"`
	Country    string `json:"country"`
	ReturnURL  string `json:"return_url"`
	RefreshURL string `json:"refresh_url"`
}

// StartFreelancerOnboardingResult reports the Stripe onboarding URL and
// connected-account state.
type StartFreelancerOnboardingResult struct {
	UserID           string `json:"user_id"`
	StripeAccountID  string `json:"stripe_account_id"`
	OnboardingURL    string `json:"onboarding_url"`
	Status           string `json:"status"`
	DetailsSubmitted bool   `json:"details_submitted"`
	ChargesEnabled   bool   `json:"charges_enabled"`
	PayoutsEnabled   bool   `json:"payouts_enabled"`
	DisabledReason   string `json:"disabled_reason"`
}
