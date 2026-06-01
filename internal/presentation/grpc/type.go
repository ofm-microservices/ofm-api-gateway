package grpc

import (
	"api-gateway/config"
	gateway "api-gateway/internal/domain"
	"context"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	authv1 "github.com/ofm-microservices/ofm-common/proto/auth/v1"
	gigv1 "github.com/ofm-microservices/ofm-common/proto/gig/v1"
	ordercheckoutv1 "github.com/ofm-microservices/ofm-common/proto/ordercheckout/v1"
	paymentconnectv1 "github.com/ofm-microservices/ofm-common/proto/paymentconnect/v1"
	registrationv1 "github.com/ofm-microservices/ofm-common/proto/registration/v1"
	reviewv1 "github.com/ofm-microservices/ofm-common/proto/review/v1"
	searchv1 "github.com/ofm-microservices/ofm-common/proto/search/v1"
	userv1 "github.com/ofm-microservices/ofm-common/proto/user/v1"
)

// Logger aliases the shared logger contract used by the gRPC adapter.
type Logger = logging.Logger

// SignUpRequest aliases the gateway-domain signup request transported over gRPC.
type SignUpRequest = gateway.SignUpRequest

// SignUpResult aliases the gateway-domain signup result returned from the saga.
type SignUpResult = gateway.SignUpResult

// VerifyEmailRequest aliases the gateway-domain verification command.
type VerifyEmailRequest = gateway.VerifyEmailRequest

// VerifyEmailResult aliases the gateway-domain verification result.
type VerifyEmailResult = gateway.VerifyEmailResult

// AuthTokensResult aliases the gateway-domain token result used by both auth
// flows.
type AuthTokensResult = gateway.AuthTokensResult

// CompleteRegistrationResult is retained as a compatibility alias for older
// registration-only code paths.
type CompleteRegistrationResult = gateway.CompleteRegistrationResult

// Client is the gateway-facing gRPC adapter for registration-saga-service.
type Client interface {
	StartRegistration(ctx context.Context, req SignUpRequest) (*SignUpResult, error)
	VerifyEmail(ctx context.Context, req VerifyEmailRequest) (*VerifyEmailResult, error)
	GetRegistrationStatus(ctx context.Context, sessionID, clientID string) (*gateway.RegistrationStatus, error)
	Close() error
}

// AuthClient is the gateway-facing gRPC adapter for auth-service token issuing.
type AuthClient interface {
	IssueRegistrationTokens(ctx context.Context, userID string) (*CompleteRegistrationResult, error)
	Close() error
}

// AuthSessionClient is the gateway-facing gRPC adapter for sign-in.
type AuthSessionClient interface {
	SignIn(ctx context.Context, req gateway.SignInRequest) (*AuthTokensResult, error)
	Refresh(ctx context.Context, req gateway.RefreshTokensRequest) (*AuthTokensResult, error)
	SignOut(ctx context.Context, req gateway.SignOutRequest) error
	Close() error
}

// GigClient is the gateway-facing gRPC adapter for gig-service draft and
// publish workflows.
type GigClient interface {
	CreateDraft(ctx context.Context, req gateway.CreateGigDraftRequest) (*gateway.Gig, error)
	UpdateBasicInfo(ctx context.Context, req gateway.UpdateGigBasicInfoRequest) (*gateway.Gig, error)
	ReplacePackages(ctx context.Context, req gateway.ReplaceGigPackagesRequest) (*gateway.Gig, error)
	ReplaceQuestions(ctx context.Context, req gateway.ReplaceGigQuestionsRequest) (*gateway.Gig, error)
	ReplaceMedia(ctx context.Context, req gateway.ReplaceGigMediaRequest) (*gateway.Gig, error)
	GetDraft(ctx context.Context, req gateway.GetGigDraftRequest) (*gateway.Gig, error)
	GetBySlug(ctx context.Context, req gateway.GetGigBySlugRequest) (*gateway.Gig, error)
	Publish(ctx context.Context, req gateway.PublishGigRequest) (*gateway.Gig, error)
	Close() error
}

// PaymentOnboardingClient is the gateway-facing gRPC adapter for freelancer
// Stripe onboarding.
type PaymentOnboardingClient interface {
	StartFreelancerOnboarding(ctx context.Context, req gateway.StartFreelancerOnboardingRequest) (*gateway.StartFreelancerOnboardingResult, error)
	Close() error
}

// OrderCheckoutClient is the gateway-facing gRPC adapter for the hybrid order
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

// ReviewClient is the gateway-facing gRPC adapter for review-service.
type ReviewClient interface {
	CreateReview(ctx context.Context, req gateway.CreateReviewRequest) (*gateway.CreateReviewResult, error)
	GetGigReviews(ctx context.Context, req gateway.GetGigReviewsRequest) (*gateway.GetGigReviewsResult, error)
	GetGigReviewsSummary(ctx context.Context, req gateway.GetGigReviewsSummaryRequest) (*gateway.ReviewSummary, error)
	GetUserRatingSummaryByUsername(ctx context.Context, req gateway.GetUserRatingSummaryByUsernameRequest) (*gateway.ReviewSummary, error)
	Close() error
}

// UserClient is the gateway-facing gRPC adapter for user-service.
type UserClient interface {
	GetDetailedUserByUsername(ctx context.Context, username string) (*gateway.User, error)
	Close() error
}

// SearchClient is the gateway-facing gRPC adapter for search-service.
type SearchClient interface {
	Search(ctx context.Context, req gateway.SearchRequest) (*gateway.SearchResponse, error)
	Close() error
}

// RegistrationMapper translates between gateway-domain signup types and the
// shared registration gRPC contract.
type RegistrationMapper interface {
	ToStartRegistrationRequest(req SignUpRequest) *registrationv1.StartRegistrationRequest
	ToSignUpResult(res *registrationv1.StartRegistrationResponse) *SignUpResult
	ToVerifyEmailRequest(req VerifyEmailRequest) *registrationv1.VerifyEmailRequest
	ToVerifyEmailResult(res *registrationv1.VerifyEmailResponse) *VerifyEmailResult
	ToRegistrationStatus(res *registrationv1.GetRegistrationStatusResponse) *gateway.RegistrationStatus
	ToCompleteRegistrationResult(res *authv1.IssueRegistrationTokensResponse) *AuthTokensResult
	ToStartRegistrationError(err error) error
	ToRegistrationStatusError(err error) error
}

// AuthSessionMapper translates between gateway sign-in types and the shared
// auth-session gRPC contract.
type AuthSessionMapper interface {
	ToSignInRequest(req gateway.SignInRequest) *authv1.SignInRequest
	ToRefreshRequest(req gateway.RefreshTokensRequest) *authv1.RefreshRequest
	ToSignOutRequest(req gateway.SignOutRequest) *authv1.SignOutRequest
	ToSignInResult(res *authv1.SignInResponse) *AuthTokensResult
	ToRefreshResult(res *authv1.RefreshResponse) *AuthTokensResult
	ToError(err error) error
	ToRefreshError(err error) error
	ToSignOutError(err error) error
}

// GigMapper translates between gateway gig types and the shared gig gRPC
// contract.
type GigMapper interface {
	ToCreateDraftRequest(req gateway.CreateGigDraftRequest) *gigv1.CreateDraftRequest
	ToCreateDraftResponse(res *gigv1.CreateDraftResponse) *gateway.Gig
	ToUpdateBasicInfoRequest(req gateway.UpdateGigBasicInfoRequest) *gigv1.UpdateBasicInfoRequest
	ToUpdateBasicInfoResponse(res *gigv1.UpdateBasicInfoResponse) *gateway.Gig
	ToReplacePackagesRequest(req gateway.ReplaceGigPackagesRequest) *gigv1.ReplacePackagesRequest
	ToReplacePackagesResponse(res *gigv1.ReplacePackagesResponse) *gateway.Gig
	ToReplaceQuestionsRequest(req gateway.ReplaceGigQuestionsRequest) *gigv1.ReplaceQuestionsRequest
	ToReplaceQuestionsResponse(res *gigv1.ReplaceQuestionsResponse) *gateway.Gig
	ToReplaceMediaRequest(req gateway.ReplaceGigMediaRequest) *gigv1.ReplaceMediaRequest
	ToReplaceMediaResponse(res *gigv1.ReplaceMediaResponse) *gateway.Gig
	ToGetDraftRequest(req gateway.GetGigDraftRequest) *gigv1.GetDraftRequest
	ToGetDraftResponse(res *gigv1.GetDraftResponse) *gateway.Gig
	ToGetBySlugRequest(req gateway.GetGigBySlugRequest) *gigv1.GetGigBySlugRequest
	ToGetBySlugResponse(res *gigv1.GetGigBySlugResponse) *gateway.Gig
	ToPublishRequest(req gateway.PublishGigRequest) *gigv1.PublishRequest
	ToPublishResponse(res *gigv1.PublishResponse) *gateway.Gig
	ToError(err error) error
}

// PaymentOnboardingMapper translates between gateway onboarding types and the
// shared payment onboarding gRPC contract.
type PaymentOnboardingMapper interface {
	ToStartFreelancerOnboardingRequest(req gateway.StartFreelancerOnboardingRequest) *paymentconnectv1.StartFreelancerOnboardingRequest
	ToStartFreelancerOnboardingResponse(res *paymentconnectv1.StartFreelancerOnboardingResponse) *gateway.StartFreelancerOnboardingResult
	ToError(err error) error
}

// OrderCheckoutMapper translates between gateway order types and the shared
// order checkout gRPC contract.
type OrderCheckoutMapper interface {
	ToStartOrderRequest(req gateway.CreateOrderRequest) *ordercheckoutv1.StartOrderRequest
	ToStartOrderResponse(res *ordercheckoutv1.StartOrderResponse) *gateway.CreateOrderResult
	ToConfirmOrderRequest(req gateway.ConfirmOrderRequest) *ordercheckoutv1.ConfirmOrderRequest
	ToConfirmOrderResponse(res *ordercheckoutv1.ConfirmOrderResponse) *gateway.ConfirmOrderResult
	ToSubmitRequirementsRequest(req gateway.SubmitOrderRequirementsRequest) *ordercheckoutv1.SubmitRequirementsRequest
	ToSubmitRequirementsResponse(res *ordercheckoutv1.SubmitRequirementsResponse) *gateway.SubmitOrderRequirementsResult
	ToSubmitMessageRequest(req gateway.SubmitOrderMessageRequest) *ordercheckoutv1.SubmitMessageRequest
	ToSubmitMessageResponse(res *ordercheckoutv1.SubmitMessageResponse) *gateway.SubmitOrderMessageResult
	ToCreateAttachmentUploadURLRequest(req gateway.CreateOrderAttachmentUploadURLRequest) *ordercheckoutv1.CreateAttachmentUploadURLRequest
	ToCreateAttachmentUploadURLResponse(res *ordercheckoutv1.CreateAttachmentUploadURLResponse) *gateway.CreateOrderAttachmentUploadURLResult
	ToCompleteAttachmentUploadRequest(req gateway.CompleteOrderAttachmentUploadRequest) *ordercheckoutv1.CompleteAttachmentUploadRequest
	ToCompleteAttachmentUploadResponse(res *ordercheckoutv1.CompleteAttachmentUploadResponse) *gateway.CompleteOrderAttachmentUploadResult
	ToDeliverOrderRequest(req gateway.DeliverOrderRequest) *ordercheckoutv1.DeliverOrderRequest
	ToDeliverOrderResponse(res *ordercheckoutv1.DeliverOrderResponse) *gateway.DeliverOrderResult
	ToAcceptDeliveryRequest(req gateway.AcceptDeliveryRequest) *ordercheckoutv1.AcceptDeliveryRequest
	ToAcceptDeliveryResponse(res *ordercheckoutv1.AcceptDeliveryResponse) *gateway.AcceptDeliveryResult
	ToRequestRevisionRequest(req gateway.RequestRevisionRequest) *ordercheckoutv1.RequestRevisionRequest
	ToRequestRevisionResponse(res *ordercheckoutv1.RequestRevisionResponse) *gateway.RequestRevisionResult
	ToOpenDisputeRequest(req gateway.OpenDisputeRequest) *ordercheckoutv1.OpenDisputeRequest
	ToOpenDisputeResponse(res *ordercheckoutv1.OpenDisputeResponse) *gateway.OpenDisputeResult
	ToError(err error) error
}

// ReviewMapper translates between gateway review types and the shared review
// gRPC contract.
type ReviewMapper interface {
	ToCreateReviewRequest(req gateway.CreateReviewRequest, buyerID string) *reviewv1.CreateReviewRequest
	ToCreateReviewResponse(res *reviewv1.CreateReviewResponse) *gateway.CreateReviewResult
	ToGetGigReviewsRequest(req gateway.GetGigReviewsRequest) *reviewv1.ListGigReviewsRequest
	ToGetGigReviewsResponse(res *reviewv1.ListGigReviewsResponse) *gateway.GetGigReviewsResult
	ToGetGigReviewsSummaryRequest(req gateway.GetGigReviewsSummaryRequest) *reviewv1.GetGigRatingSummaryRequest
	ToGetGigReviewsSummaryResponse(res *reviewv1.RatingSummary) *gateway.ReviewSummary
	ToGetUserRatingSummaryByUsernameRequest(req gateway.GetUserRatingSummaryByUsernameRequest) *reviewv1.GetUserRatingSummaryByUsernameRequest
	ToGetUserRatingSummaryByUsernameResponse(res *reviewv1.RatingSummary) *gateway.ReviewSummary
	ToError(err error) error
}

// SearchMapper translates between gateway search types and the shared search
// gRPC contract.
type SearchMapper interface {
	ToSearchRequest(req gateway.SearchRequest) *searchv1.SearchRequest
	ToSearchResponse(res *searchv1.SearchResponse) *gateway.SearchResponse
}

// UserMapper translates between gateway user types and the shared user gRPC
// contract.
type UserMapper interface {
	ToGetDetailedUserRequest(username string) *userv1.GetDetailedUserByUsernameRequest
	ToGetDetailedUserResponse(res *userv1.GetDetailedUserByUsernameResponse) *gateway.User
	ToError(err error) error
}

// UserServiceConfig aliases the outbound user-service gRPC client configuration.
type UserServiceConfig = config.UserServiceConfig

// RegistrationSagaConfig aliases the outbound saga gRPC client configuration.
type RegistrationSagaConfig = config.RegistrationSagaConfig

// AuthServiceConfig aliases the outbound auth gRPC client configuration.
type AuthServiceConfig = config.AuthServiceConfig

// GigServiceConfig aliases the outbound gig gRPC client configuration.
type GigServiceConfig = config.GigServiceConfig

// PaymentServiceConfig aliases the outbound payment gRPC client configuration.
type PaymentServiceConfig = config.PaymentServiceConfig

// OrderSagaConfig aliases the outbound order-saga gRPC client configuration.
type OrderSagaConfig = config.OrderSagaConfig

// ReviewServiceConfig aliases the outbound review-service gRPC client configuration.
type ReviewServiceConfig = config.ReviewServiceConfig

// SearchServiceConfig aliases the outbound search-service gRPC client configuration.
type SearchServiceConfig = config.SearchServiceConfig
