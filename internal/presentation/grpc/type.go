package grpc

import (
	"api-gateway/config"
	gateway "api-gateway/internal/domain"
	"context"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	authv1 "github.com/ofm-microservices/ofm-common/proto/auth/v1"
	gigv1 "github.com/ofm-microservices/ofm-common/proto/gig/v1"
	registrationv1 "github.com/ofm-microservices/ofm-common/proto/registration/v1"
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

// CompleteRegistrationResult aliases the gateway-domain token result.
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

// GigClient is the gateway-facing gRPC adapter for gig-service draft and
// publish workflows.
type GigClient interface {
	CreateDraft(ctx context.Context, req gateway.CreateGigDraftRequest) (*gateway.Gig, error)
	UpdateBasicInfo(ctx context.Context, req gateway.UpdateGigBasicInfoRequest) (*gateway.Gig, error)
	ReplacePackages(ctx context.Context, req gateway.ReplaceGigPackagesRequest) (*gateway.Gig, error)
	ReplaceQuestions(ctx context.Context, req gateway.ReplaceGigQuestionsRequest) (*gateway.Gig, error)
	ReplaceMedia(ctx context.Context, req gateway.ReplaceGigMediaRequest) (*gateway.Gig, error)
	GetDraft(ctx context.Context, req gateway.GetGigDraftRequest) (*gateway.Gig, error)
	Publish(ctx context.Context, req gateway.PublishGigRequest) (*gateway.Gig, error)
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
	ToCompleteRegistrationResult(res *authv1.IssueRegistrationTokensResponse) *CompleteRegistrationResult
	ToStartRegistrationError(err error) error
	ToRegistrationStatusError(err error) error
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
	ToPublishRequest(req gateway.PublishGigRequest) *gigv1.PublishRequest
	ToPublishResponse(res *gigv1.PublishResponse) *gateway.Gig
	ToError(err error) error
}

// RegistrationSagaConfig aliases the outbound saga gRPC client configuration.
type RegistrationSagaConfig = config.RegistrationSagaConfig

// AuthServiceConfig aliases the outbound auth gRPC client configuration.
type AuthServiceConfig = config.AuthServiceConfig

// GigServiceConfig aliases the outbound gig gRPC client configuration.
type GigServiceConfig = config.GigServiceConfig
