package service

import (
	gateway "api-gateway/internal/domain"
	"context"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
)

// RegistrationPublisher is the outbound boundary used to start registration in
// another service.
type RegistrationPublisher interface {
	// StartRegistration delegates registration start to the saga service.
	StartRegistration(ctx context.Context, req gateway.SignUpRequest) (*gateway.SignUpResult, error)
	// VerifyEmail delegates email verification to the saga service.
	VerifyEmail(ctx context.Context, req gateway.VerifyEmailRequest) (*gateway.VerifyEmailResult, error)
	// GetRegistrationStatus returns saga state before token exchange.
	GetRegistrationStatus(ctx context.Context, sessionID, clientID string) (*gateway.RegistrationStatus, error)
}

// TokenIssuer is the outbound auth-service boundary for final token exchange.
type TokenIssuer interface {
	// IssueRegistrationTokens creates login tokens after saga completion.
	IssueRegistrationTokens(ctx context.Context, userID string) (*gateway.AuthTokensResult, error)
	Close() error
}

// AuthSessionClient is the outbound auth-service boundary used for sign-in.
type AuthSessionClient interface {
	SignIn(ctx context.Context, req gateway.SignInRequest) (*gateway.AuthTokensResult, error)
	Refresh(ctx context.Context, req gateway.RefreshTokensRequest) (*gateway.AuthTokensResult, error)
	SignOut(ctx context.Context, req gateway.SignOutRequest) error
	Close() error
}

// RegistrationService validates the public signup request and delegates
// orchestration to registration-saga-service.
type RegistrationService interface {
	// SignUp validates the public request and starts the registration saga.
	SignUp(ctx context.Context, req gateway.SignUpRequest) (*gateway.SignUpResult, error)
	// VerifyEmail accepts the emailed code and starts saga completion.
	VerifyEmail(ctx context.Context, req gateway.VerifyEmailRequest) (*gateway.VerifyEmailResult, error)
	// CompleteRegistration exchanges a completed saga for auth-owned tokens.
	CompleteRegistration(ctx context.Context, req gateway.CompleteRegistrationRequest) (*gateway.AuthTokensResult, error)
}

// AuthSessionService validates public sign-in requests and delegates to auth-service.
type AuthSessionService interface {
	SignIn(ctx context.Context, req gateway.SignInRequest) (*gateway.AuthTokensResult, error)
	Refresh(ctx context.Context, req gateway.RefreshTokensRequest) (*gateway.AuthTokensResult, error)
	SignOut(ctx context.Context, req gateway.SignOutRequest) error
}

// GigService validates the public gig draft workflow and delegates
// orchestration to gig-service.
type GigService interface {
	CreateDraft(ctx context.Context, req gateway.CreateGigDraftRequest) (*gateway.Gig, error)
	UpdateBasicInfo(ctx context.Context, req gateway.UpdateGigBasicInfoRequest) (*gateway.Gig, error)
	ReplacePackages(ctx context.Context, req gateway.ReplaceGigPackagesRequest) (*gateway.Gig, error)
	ReplaceQuestions(ctx context.Context, req gateway.ReplaceGigQuestionsRequest) (*gateway.Gig, error)
	ReplaceMedia(ctx context.Context, req gateway.ReplaceGigMediaRequest) (*gateway.Gig, error)
	GetDraft(ctx context.Context, req gateway.GetGigDraftRequest) (*gateway.Gig, error)
	GetBySlug(ctx context.Context, req gateway.GetGigBySlugRequest) (*gateway.Gig, error)
	Publish(ctx context.Context, req gateway.PublishGigRequest) (*gateway.Gig, error)
}

// UserClient is the outbound gRPC boundary for user-service.
type UserClient interface {
	GetDetailedUserByUsername(ctx context.Context, username string) (*gateway.User, error)
	Close() error
}

// GigPublisher is the outbound boundary used to manage gig drafts in another
// service.
type GigPublisher interface {
	CreateDraft(ctx context.Context, req gateway.CreateGigDraftRequest) (*gateway.Gig, error)
	UpdateBasicInfo(ctx context.Context, req gateway.UpdateGigBasicInfoRequest) (*gateway.Gig, error)
	ReplacePackages(ctx context.Context, req gateway.ReplaceGigPackagesRequest) (*gateway.Gig, error)
	ReplaceQuestions(ctx context.Context, req gateway.ReplaceGigQuestionsRequest) (*gateway.Gig, error)
	ReplaceMedia(ctx context.Context, req gateway.ReplaceGigMediaRequest) (*gateway.Gig, error)
	GetDraft(ctx context.Context, req gateway.GetGigDraftRequest) (*gateway.Gig, error)
	GetBySlug(ctx context.Context, req gateway.GetGigBySlugRequest) (*gateway.Gig, error)
	Publish(ctx context.Context, req gateway.PublishGigRequest) (*gateway.Gig, error)
}

// OrderPublisher is the outbound boundary used to start order sagas in the
// order-saga service.
type OrderPublisher interface {
	StartOrder(ctx context.Context, req gateway.CreateOrderRequest) (*gateway.CreateOrderResult, error)
}

// OrderCheckoutClient is the outbound gRPC boundary for the hybrid order
// checkout flow.
type OrderCheckoutClient interface {
	StartOrder(ctx context.Context, req gateway.CreateOrderRequest) (*gateway.CreateOrderResult, error)
	ConfirmOrder(ctx context.Context, req gateway.ConfirmOrderRequest) (*gateway.ConfirmOrderResult, error)
	SubmitRequirements(ctx context.Context, req gateway.SubmitOrderRequirementsRequest) (*gateway.SubmitOrderRequirementsResult, error)
	SubmitMessage(ctx context.Context, req gateway.SubmitOrderMessageRequest) (*gateway.SubmitOrderMessageResult, error)
	CreateAttachmentUploadURL(ctx context.Context, req gateway.CreateOrderAttachmentUploadURLRequest) (*gateway.CreateOrderAttachmentUploadURLResult, error)
	CompleteAttachmentUpload(ctx context.Context, req gateway.CompleteOrderAttachmentUploadRequest) (*gateway.CompleteOrderAttachmentUploadResult, error)
	DeliverOrder(ctx context.Context, req gateway.DeliverOrderRequest) (*gateway.DeliverOrderResult, error)
	AcceptDelivery(ctx context.Context, req gateway.AcceptDeliveryRequest) (*gateway.AcceptDeliveryResult, error)
	RequestRevision(ctx context.Context, req gateway.RequestRevisionRequest) (*gateway.RequestRevisionResult, error)
	OpenDispute(ctx context.Context, req gateway.OpenDisputeRequest) (*gateway.OpenDisputeResult, error)
	Close() error
}

// ReviewClient is the outbound gRPC boundary for review-service.
type ReviewClient interface {
	CreateReview(ctx context.Context, req gateway.CreateReviewRequest) (*gateway.CreateReviewResult, error)
	GetGigReviews(ctx context.Context, req gateway.GetGigReviewsRequest) (*gateway.GetGigReviewsResult, error)
	GetGigReviewsSummary(ctx context.Context, req gateway.GetGigReviewsSummaryRequest) (*gateway.ReviewSummary, error)
	GetUserRatingSummaryByUsername(ctx context.Context, req gateway.GetUserRatingSummaryByUsernameRequest) (*gateway.ReviewSummary, error)
	Close() error
}

// SearchClient is the outbound gRPC boundary for search-service.
type SearchClient interface {
	Search(ctx context.Context, req gateway.SearchRequest) (*gateway.SearchResponse, error)
	Close() error
}

// OrderService validates public create-order requests and delegates to the
// saga boundary.
type OrderService interface {
	CreateOrder(ctx context.Context, req gateway.CreateOrderRequest) (*gateway.CreateOrderResult, error)
	ConfirmOrder(ctx context.Context, req gateway.ConfirmOrderRequest) (*gateway.ConfirmOrderResult, error)
	SubmitRequirements(ctx context.Context, req gateway.SubmitOrderRequirementsRequest) (*gateway.SubmitOrderRequirementsResult, error)
	SubmitMessage(ctx context.Context, req gateway.SubmitOrderMessageRequest) (*gateway.SubmitOrderMessageResult, error)
	CreateAttachmentUploadURL(ctx context.Context, req gateway.CreateOrderAttachmentUploadURLRequest) (*gateway.CreateOrderAttachmentUploadURLResult, error)
	CompleteAttachmentUpload(ctx context.Context, req gateway.CompleteOrderAttachmentUploadRequest) (*gateway.CompleteOrderAttachmentUploadResult, error)
	DeliverOrder(ctx context.Context, req gateway.DeliverOrderRequest) (*gateway.DeliverOrderResult, error)
	AcceptDelivery(ctx context.Context, req gateway.AcceptDeliveryRequest) (*gateway.AcceptDeliveryResult, error)
	RequestRevision(ctx context.Context, req gateway.RequestRevisionRequest) (*gateway.RequestRevisionResult, error)
	OpenDispute(ctx context.Context, req gateway.OpenDisputeRequest) (*gateway.OpenDisputeResult, error)
}

// ReviewService validates public review requests and delegates to review-service.
type ReviewService interface {
	CreateReview(ctx context.Context, req gateway.CreateReviewRequest) (*gateway.CreateReviewResult, error)
}

// SearchService validates public search requests and delegates to search-service.
type SearchService interface {
	Search(ctx context.Context, req gateway.SearchRequest) (*gateway.SearchResponse, error)
}

// PaymentOnboardingService validates public freelancer onboarding requests and
// delegates to payment-service.
type PaymentOnboardingService interface {
	StartFreelancerOnboarding(ctx context.Context, req gateway.StartFreelancerOnboardingRequest) (*gateway.StartFreelancerOnboardingResult, error)
}

// PublicGigService resolves the public gig detail page with review enrichment.
type PublicGigService interface {
	GetBySlug(ctx context.Context, req gateway.GetGigBySlugRequest) (*gateway.Gig, error)
}

// PaymentOnboardingPublisher is the outbound boundary used to start
// freelancer onboarding in payment-service.
type PaymentOnboardingPublisher interface {
	StartFreelancerOnboarding(ctx context.Context, req gateway.StartFreelancerOnboardingRequest) (*gateway.StartFreelancerOnboardingResult, error)
}

// Logger aliases the shared structured logger used by the application layer.
type Logger = logging.Logger
