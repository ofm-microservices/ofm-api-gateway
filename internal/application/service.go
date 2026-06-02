package service

import (
	gateway "api-gateway/internal/domain"
	"context"
	"github.com/google/uuid"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"net/mail"
	"strings"
	"time"
)

type registrationService struct {
	client RegistrationPublisher
	tokens TokenIssuer
	log    Logger
}

type authSessionService struct {
	client AuthSessionClient
	log    Logger
}

type authMeService struct {
	client UserClient
	log    Logger
}

type orderService struct {
	client OrderCheckoutClient
	log    Logger
}

type onboardingService struct {
	client PaymentOnboardingPublisher
	log    Logger
}

type reviewService struct {
	client ReviewClient
	log    Logger
}

type searchService struct {
	client SearchClient
	log    Logger
}

// New constructs the application service responsible for starting
// registrations through the saga boundary.
func New(client RegistrationPublisher, tokens TokenIssuer, log Logger) (RegistrationService, error) {
	if client == nil {
		return nil, ErrNilRegistrationClient
	}
	if tokens == nil {
		return nil, ErrNilTokenIssuer
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	return &registrationService{
		client: client,
		tokens: tokens,
		log:    log.With(logging.String("module", "application")),
	}, nil
}

// NewAuthSession constructs the application service responsible for signing
// into auth-service through the session boundary.
func NewAuthSession(client AuthSessionClient, log Logger) (AuthSessionService, error) {
	if client == nil {
		return nil, ErrNilAuthSessionClient
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	return &authSessionService{
		client: client,
		log:    log.With(logging.String("module", "auth-session-application")),
	}, nil
}

// NewAuthMe constructs the application service responsible for resolving the
// current authenticated user preview from user-service.
func NewAuthMe(client UserClient, log Logger) (AuthMeService, error) {
	if client == nil {
		return nil, ErrNilUserClient
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	return &authMeService{
		client: client,
		log:    log.With(logging.String("module", "auth-me-application")),
	}, nil
}

// NewOrder constructs the application service responsible for order checkout
// orchestration through the saga boundary.
func NewOrder(client OrderCheckoutClient, log Logger) (OrderService, error) {
	if client == nil {
		return nil, ErrNilOrderClient
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	return &orderService{
		client: client,
		log:    log.With(logging.String("module", "order-application")),
	}, nil
}

// NewPaymentOnboarding constructs the service responsible for starting Stripe
// Connect onboarding through payment-service.
func NewPaymentOnboarding(client PaymentOnboardingPublisher, log Logger) (PaymentOnboardingService, error) {
	if client == nil {
		return nil, ErrNilPaymentOnboardingClient
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	return &onboardingService{
		client: client,
		log:    log.With(logging.String("module", "onboarding-application")),
	}, nil
}

// NewReview constructs the application service responsible for buyer review submissions.
func NewReview(client ReviewClient, log Logger) (ReviewService, error) {
	if client == nil {
		return nil, ErrNilReviewClient
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	return &reviewService{
		client: client,
		log:    log.With(logging.String("module", "review-application")),
	}, nil
}

func (s *registrationService) SignUp(ctx context.Context, req gateway.SignUpRequest) (*gateway.SignUpResult, error) {
	started := time.Now()
	log := logging.WithContext(ctx, s.log)
	log.Info("sign up request received",
		logging.Operation("registration.sign_up"),
		logging.String("email", req.Email),
		logging.String("username", req.Username),
	)

	if _, err := mail.ParseAddress(strings.TrimSpace(req.Email)); err != nil {
		return nil, gateway.ErrInvalidEmail
	}
	if strings.TrimSpace(req.Username) == "" {
		return nil, gateway.ErrInvalidUsername
	}
	if len(req.Password) < 8 {
		return nil, gateway.ErrInvalidPassword
	}

	result, err := s.client.StartRegistration(ctx, gateway.SignUpRequest{
		ClientID:  strings.TrimSpace(req.ClientID),
		Email:     strings.TrimSpace(req.Email),
		Password:  req.Password,
		Username:  strings.TrimSpace(req.Username),
		FirstName: strings.TrimSpace(req.FirstName),
		Surname:   strings.TrimSpace(req.Surname),
	})
	if err != nil {
		log.Error("failed to start registration",
			logging.Operation("registration.sign_up"),
			logging.Attempt(1),
			logging.Retryable(false),
			logging.DurationMS(time.Since(started)),
			logging.String("email", req.Email),
			logging.Err(err),
		)
		return nil, err
	}
	if result.Status == "conflict" {
		return nil, &gateway.RegistrationConflictError{
			State:         result.ConflictState,
			UsernameTaken: result.UsernameTaken,
			EmailTaken:    result.EmailTaken,
		}
	}

	log.Info("registration started",
		logging.Operation("registration.sign_up"),
		logging.DurationMS(time.Since(started)),
		logging.String("email", req.Email),
		logging.String("session_id", result.SessionID),
	)
	return result, nil
}

// NewSearch constructs the application service responsible for public gig search.
func NewSearch(client SearchClient, log Logger) (SearchService, error) {
	if client == nil {
		return nil, ErrNilSearchClient
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	return &searchService{
		client: client,
		log:    log.With(logging.String("module", "search-application")),
	}, nil
}

func (s *registrationService) VerifyEmail(ctx context.Context, req gateway.VerifyEmailRequest) (*gateway.VerifyEmailResult, error) {
	log := logging.WithContext(ctx, s.log)
	sessionID := strings.TrimSpace(req.SessionID)
	clientID := strings.TrimSpace(req.ClientID)
	code := strings.TrimSpace(req.Code)
	if sessionID == "" {
		return nil, gateway.ErrInvalidSessionID
	}
	if clientID == "" {
		return nil, gateway.ErrInvalidClientID
	}
	if code == "" {
		return nil, gateway.ErrInvalidVerificationCode
	}

	result, err := s.client.VerifyEmail(ctx, gateway.VerifyEmailRequest{
		SessionID: sessionID,
		ClientID:  clientID,
		Code:      code,
	})
	if err != nil {
		log.Error("failed to verify registration email",
			logging.Operation("registration.verify_email"),
			logging.Attempt(1),
			logging.Retryable(false),
			logging.String("session_id", sessionID),
			logging.Err(err),
		)
		return nil, err
	}

	return result, nil
}

func (s *registrationService) CompleteRegistration(ctx context.Context, req gateway.CompleteRegistrationRequest) (*gateway.AuthTokensResult, error) {
	log := logging.WithContext(ctx, s.log)
	sessionID := strings.TrimSpace(req.SessionID)
	clientID := strings.TrimSpace(req.ClientID)
	if sessionID == "" {
		return nil, gateway.ErrInvalidSessionID
	}
	if clientID == "" {
		return nil, gateway.ErrInvalidClientID
	}

	status, err := s.client.GetRegistrationStatus(ctx, sessionID, clientID)
	if err != nil {
		log.Error("failed to get registration status",
			logging.Operation("registration.status"),
			logging.Attempt(1),
			logging.Retryable(false),
			logging.String("session_id", sessionID),
			logging.Err(err),
		)
		return nil, err
	}
	if status.Status == "tokens_claimed" {
		return nil, gateway.ErrRegistrationAlreadyClaimed
	}
	if status.Status != "completed" {
		return nil, gateway.ErrRegistrationNotCompleted
	}

	result, err := s.tokens.IssueRegistrationTokens(ctx, status.UserID)
	if err != nil {
		log.Error("failed to issue registration tokens",
			logging.Operation("registration.issue_tokens"),
			logging.Attempt(1),
			logging.Retryable(false),
			logging.String("user_id", status.UserID),
			logging.Err(err),
		)
		return nil, err
	}

	return result, nil
}

func (s *authSessionService) SignIn(ctx context.Context, req gateway.SignInRequest) (*gateway.AuthTokensResult, error) {
	log := logging.WithContext(ctx, s.log)
	identifier := strings.TrimSpace(req.Identifier)
	password := req.Password
	if identifier == "" {
		return nil, gateway.ErrInvalidIdentifier
	}
	if strings.TrimSpace(password) == "" {
		return nil, gateway.ErrInvalidPassword
	}

	result, err := s.client.SignIn(ctx, gateway.SignInRequest{
		Identifier: identifier,
		Password:   password,
	})
	if err != nil {
		log.Error("failed to sign in",
			logging.Operation("auth.sign_in"),
			logging.Attempt(1),
			logging.Retryable(false),
			logging.String("identifier", identifier),
			logging.Err(err),
		)
		return nil, err
	}

	return result, nil
}

func (s *authSessionService) Refresh(ctx context.Context, req gateway.RefreshTokensRequest) (*gateway.AuthTokensResult, error) {
	log := logging.WithContext(ctx, s.log)
	refreshToken := strings.TrimSpace(req.RefreshToken)
	if refreshToken == "" {
		return nil, gateway.ErrInvalidRefreshToken
	}

	result, err := s.client.Refresh(ctx, gateway.RefreshTokensRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		log.Error("failed to refresh tokens",
			logging.Operation("auth.refresh"),
			logging.Attempt(1),
			logging.Retryable(false),
			logging.Err(err),
		)
		return nil, err
	}

	return result, nil
}

func (s *authSessionService) SignOut(ctx context.Context, req gateway.SignOutRequest) error {
	log := logging.WithContext(ctx, s.log)
	refreshToken := strings.TrimSpace(req.RefreshToken)
	if refreshToken == "" {
		return gateway.ErrInvalidRefreshToken
	}

	if err := s.client.SignOut(ctx, gateway.SignOutRequest{
		RefreshToken: refreshToken,
	}); err != nil {
		log.Error("failed to sign out",
			logging.Operation("auth.sign_out"),
			logging.Attempt(1),
			logging.Retryable(false),
			logging.Err(err),
		)
		return err
	}

	return nil
}

func (s *orderService) CreateOrder(ctx context.Context, req gateway.CreateOrderRequest) (*gateway.CreateOrderResult, error) {
	log := logging.WithContext(ctx, s.log)
	buyerID := strings.TrimSpace(req.BuyerID)
	if buyerID == "" {
		return nil, gateway.ErrInvalidOrderBuyerID
	}
	buyerEmail := strings.TrimSpace(req.BuyerEmail)
	if buyerEmail != "" {
		if _, err := mail.ParseAddress(buyerEmail); err != nil {
			return nil, gateway.ErrInvalidOrderBuyerEmail
		}
	}
	if buyerEmail == "" {
		buyerEmail = ""
	}
	connectionID := strings.TrimSpace(req.RealtimeConnectionID)
	if connectionID != "" && !validRealtimeConnectionID(connectionID) {
		return nil, gateway.ErrInvalidOrderConnectionID
	}
	gigID := strings.TrimSpace(req.GigID)
	if gigID == "" {
		return nil, gateway.ErrInvalidOrderGigID
	}
	packageID := strings.TrimSpace(req.PackageID)
	if packageID == "" {
		return nil, gateway.ErrInvalidOrderPackage
	}

	result, err := s.client.StartOrder(ctx, gateway.CreateOrderRequest{
		BuyerID:              buyerID,
		BuyerEmail:           buyerEmail,
		RealtimeConnectionID: connectionID,
		GigID:                gigID,
		PackageID:            packageID,
		IdempotencyKey:       uuid.Must(uuid.NewV7()).String(),
		RequestedAt:          time.Now().UTC().Format(time.RFC3339Nano),
	})
	if err != nil {
		log.Error("failed to start order",
			logging.Operation("order.create"),
			logging.Attempt(1),
			logging.Retryable(false),
			logging.String("buyer_id", buyerID),
			logging.String("gig_id", gigID),
			logging.Err(err),
		)
		return nil, err
	}
	if result == nil {
		result = &gateway.CreateOrderResult{}
	}
	if strings.TrimSpace(result.SagaID) == "" {
		result.SagaID = uuid.Must(uuid.NewV7()).String()
	}
	if strings.TrimSpace(result.Status) == "" {
		result.Status = "requirements_pending"
	}
	return result, nil
}

func (s *orderService) ConfirmOrder(ctx context.Context, req gateway.ConfirmOrderRequest) (*gateway.ConfirmOrderResult, error) {
	log := logging.WithContext(ctx, s.log)
	orderID := strings.TrimSpace(req.OrderID)
	if orderID == "" {
		return nil, gateway.ErrInvalidOrderGigID
	}
	connectionID := strings.TrimSpace(req.RealtimeConnectionID)
	if connectionID != "" && !validRealtimeConnectionID(connectionID) {
		return nil, gateway.ErrInvalidOrderConnectionID
	}

	result, err := s.client.ConfirmOrder(ctx, gateway.ConfirmOrderRequest{
		OrderID:              orderID,
		BuyerID:              strings.TrimSpace(req.BuyerID),
		RealtimeConnectionID: connectionID,
		IdempotencyKey:       strings.TrimSpace(req.IdempotencyKey),
		RequestedAt:          time.Now().UTC().Format(time.RFC3339Nano),
	})
	if err != nil {
		log.Error("failed to confirm order",
			logging.Operation("order.confirm"),
			logging.Attempt(1),
			logging.Retryable(false),
			logging.String("order_id", orderID),
			logging.Err(err),
		)
		return nil, err
	}
	if result == nil {
		result = &gateway.ConfirmOrderResult{}
	}
	if strings.TrimSpace(result.Status) == "" {
		result.Status = "payment_pending"
	}
	return result, nil
}

func (s *orderService) SubmitRequirements(ctx context.Context, req gateway.SubmitOrderRequirementsRequest) (*gateway.SubmitOrderRequirementsResult, error) {
	return s.client.SubmitRequirements(ctx, req)
}

func (s *orderService) SubmitMessage(ctx context.Context, req gateway.SubmitOrderMessageRequest) (*gateway.SubmitOrderMessageResult, error) {
	return s.client.SubmitMessage(ctx, req)
}

func (s *orderService) CreateAttachmentUploadURL(ctx context.Context, req gateway.CreateOrderAttachmentUploadURLRequest) (*gateway.CreateOrderAttachmentUploadURLResult, error) {
	return s.client.CreateAttachmentUploadURL(ctx, req)
}

func (s *orderService) CompleteAttachmentUpload(ctx context.Context, req gateway.CompleteOrderAttachmentUploadRequest) (*gateway.CompleteOrderAttachmentUploadResult, error) {
	return s.client.CompleteAttachmentUpload(ctx, req)
}

func (s *orderService) DeliverOrder(ctx context.Context, req gateway.DeliverOrderRequest) (*gateway.DeliverOrderResult, error) {
	log := logging.WithContext(ctx, s.log)
	orderID := strings.TrimSpace(req.OrderID)
	sellerID := strings.TrimSpace(req.SellerID)
	if orderID == "" {
		return nil, gateway.ErrInvalidOrderID
	}
	if sellerID == "" {
		return nil, gateway.ErrInvalidOrderSellerID
	}
	if strings.TrimSpace(req.DeliveryMessage) == "" {
		return nil, gateway.ErrInvalidOrderDeliveryMessage
	}
	result, err := s.client.DeliverOrder(ctx, gateway.DeliverOrderRequest{
		OrderID:         orderID,
		SellerID:        sellerID,
		DeliveryMessage: strings.TrimSpace(req.DeliveryMessage),
		AttachmentIDs:   req.AttachmentIDs,
		RequestedAt:     time.Now().UTC().Format(time.RFC3339Nano),
	})
	if err != nil {
		log.Error("failed to deliver order",
			logging.Operation("order.deliver"),
			logging.Attempt(1),
			logging.Retryable(false),
			logging.String("order_id", orderID),
			logging.Err(err),
		)
		return nil, err
	}
	return result, nil
}

func (s *orderService) AcceptDelivery(ctx context.Context, req gateway.AcceptDeliveryRequest) (*gateway.AcceptDeliveryResult, error) {
	log := logging.WithContext(ctx, s.log)
	orderID := strings.TrimSpace(req.OrderID)
	buyerID := strings.TrimSpace(req.BuyerID)
	if orderID == "" {
		return nil, gateway.ErrInvalidOrderID
	}
	if buyerID == "" {
		return nil, gateway.ErrInvalidOrderBuyerID
	}
	result, err := s.client.AcceptDelivery(ctx, gateway.AcceptDeliveryRequest{
		OrderID:     orderID,
		BuyerID:     buyerID,
		RequestedAt: time.Now().UTC().Format(time.RFC3339Nano),
	})
	if err != nil {
		log.Error("failed to accept order delivery",
			logging.Operation("order.accept_delivery"),
			logging.Attempt(1),
			logging.Retryable(false),
			logging.String("order_id", orderID),
			logging.Err(err),
		)
		return nil, err
	}
	return result, nil
}

func (s *orderService) RequestRevision(ctx context.Context, req gateway.RequestRevisionRequest) (*gateway.RequestRevisionResult, error) {
	log := logging.WithContext(ctx, s.log)
	orderID := strings.TrimSpace(req.OrderID)
	buyerID := strings.TrimSpace(req.BuyerID)
	if orderID == "" {
		return nil, gateway.ErrInvalidOrderID
	}
	if buyerID == "" {
		return nil, gateway.ErrInvalidOrderBuyerID
	}
	if strings.TrimSpace(req.Reason) == "" {
		return nil, gateway.ErrInvalidOrderReason
	}
	result, err := s.client.RequestRevision(ctx, gateway.RequestRevisionRequest{
		OrderID:     orderID,
		BuyerID:     buyerID,
		Reason:      strings.TrimSpace(req.Reason),
		RequestedAt: time.Now().UTC().Format(time.RFC3339Nano),
	})
	if err != nil {
		log.Error("failed to request order revision",
			logging.Operation("order.request_revision"),
			logging.Attempt(1),
			logging.Retryable(false),
			logging.String("order_id", orderID),
			logging.Err(err),
		)
		return nil, err
	}
	return result, nil
}

func (s *orderService) OpenDispute(ctx context.Context, req gateway.OpenDisputeRequest) (*gateway.OpenDisputeResult, error) {
	log := logging.WithContext(ctx, s.log)
	orderID := strings.TrimSpace(req.OrderID)
	buyerID := strings.TrimSpace(req.BuyerID)
	if orderID == "" {
		return nil, gateway.ErrInvalidOrderID
	}
	if buyerID == "" {
		return nil, gateway.ErrInvalidOrderBuyerID
	}
	if strings.TrimSpace(req.Reason) == "" {
		return nil, gateway.ErrInvalidOrderReason
	}
	result, err := s.client.OpenDispute(ctx, gateway.OpenDisputeRequest{
		OrderID:     orderID,
		BuyerID:     buyerID,
		Reason:      strings.TrimSpace(req.Reason),
		RequestedAt: time.Now().UTC().Format(time.RFC3339Nano),
	})
	if err != nil {
		log.Error("failed to open order dispute",
			logging.Operation("order.open_dispute"),
			logging.Attempt(1),
			logging.Retryable(false),
			logging.String("order_id", orderID),
			logging.Err(err),
		)
		return nil, err
	}
	return result, nil
}

func (s *onboardingService) StartFreelancerOnboarding(ctx context.Context, req gateway.StartFreelancerOnboardingRequest) (*gateway.StartFreelancerOnboardingResult, error) {
	log := logging.WithContext(ctx, s.log)
	userID := strings.TrimSpace(req.UserID)
	if userID == "" {
		return nil, gateway.ErrInvalidFreelancerOnboarding
	}

	result, err := s.client.StartFreelancerOnboarding(ctx, gateway.StartFreelancerOnboardingRequest{
		UserID:  userID,
		Country: strings.TrimSpace(req.Country),
	})
	if err != nil {
		log.Error("failed to start freelancer onboarding",
			logging.Operation("payment.onboarding.start"),
			logging.Attempt(1),
			logging.Retryable(false),
			logging.String("user_id", userID),
			logging.Err(err),
		)
		return nil, err
	}
	return result, nil
}

func (s *reviewService) CreateReview(ctx context.Context, req gateway.CreateReviewRequest) (*gateway.CreateReviewResult, error) {
	if strings.TrimSpace(req.OrderID) == "" {
		return nil, gateway.ErrInvalidOrderID
	}
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, gateway.ErrInvalidReviewContent
	}
	if req.Rating < 1 || req.Rating > 5 {
		return nil, gateway.ErrInvalidReviewContent
	}
	if strings.TrimSpace(req.BuyerID) == "" {
		return nil, gateway.ErrInvalidOrderBuyerID
	}
	res, err := s.client.CreateReview(ctx, gateway.CreateReviewRequest{
		OrderID:     strings.TrimSpace(req.OrderID),
		BuyerID:     strings.TrimSpace(req.BuyerID),
		Content:     content,
		Rating:      req.Rating,
		RequestedAt: strings.TrimSpace(req.RequestedAt),
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *searchService) Search(ctx context.Context, req gateway.SearchRequest) (*gateway.SearchResponse, error) {
	return s.client.Search(ctx, req)
}

func (s *authMeService) GetMe(ctx context.Context, userID string) (*gateway.User, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, gateway.ErrInvalidUserID
	}

	user, err := s.client.GetUserPreviewByID(ctx, userID)
	if err != nil {
		s.log.Error("failed to resolve auth me user preview",
			logging.Operation("auth_me.get_me"),
			logging.String("user_id", userID),
			logging.Err(err),
		)
		return nil, err
	}

	return user, nil
}

func validRealtimeConnectionID(id string) bool {
	parts := strings.Split(strings.TrimSpace(id), ".")
	return len(parts) == 2 && strings.TrimSpace(parts[0]) != "" && strings.TrimSpace(parts[1]) != ""
}
