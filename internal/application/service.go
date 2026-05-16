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

type orderService struct {
	client OrderPublisher
	log    Logger
}

type onboardingService struct {
	client PaymentOnboardingPublisher
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

// NewOrder constructs the application service responsible for starting order
// sagas through the NATS boundary.
func NewOrder(client OrderPublisher, log Logger) (OrderService, error) {
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

func (s *registrationService) CompleteRegistration(ctx context.Context, req gateway.CompleteRegistrationRequest) (*gateway.CompleteRegistrationResult, error) {
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

func (s *orderService) CreateOrder(ctx context.Context, req gateway.CreateOrderRequest) (*gateway.CreateOrderResult, error) {
	log := logging.WithContext(ctx, s.log)
	buyerID := strings.TrimSpace(req.BuyerID)
	if buyerID == "" {
		return nil, gateway.ErrInvalidOrderBuyerID
	}
	buyerEmail := strings.TrimSpace(req.BuyerEmail)
	if _, err := mail.ParseAddress(buyerEmail); err != nil {
		return nil, gateway.ErrInvalidOrderBuyerEmail
	}
	connectionID := strings.TrimSpace(req.RealtimeConnectionID)
	if connectionID != "" && !validRealtimeConnectionID(connectionID) {
		return nil, gateway.ErrInvalidOrderConnectionID
	}
	gigID := strings.TrimSpace(req.GigID)
	if gigID == "" {
		return nil, gateway.ErrInvalidOrderGigID
	}
	gigTitle := strings.TrimSpace(req.GigTitle)
	if gigTitle == "" {
		return nil, gateway.ErrInvalidOrderTitle
	}
	packageID := strings.TrimSpace(req.PackageID)
	if packageID == "" {
		return nil, gateway.ErrInvalidOrderPackage
	}
	packageTier := strings.TrimSpace(req.PackageTier)
	packageDescription := strings.TrimSpace(req.PackageDescription)
	if packageTier == "" || packageDescription == "" {
		return nil, gateway.ErrInvalidOrderPackage
	}
	if req.PackageDeliveryDays <= 0 {
		return nil, gateway.ErrInvalidPackageDeliveryDays
	}
	if req.PriceCents <= 0 {
		return nil, gateway.ErrInvalidOrderPrice
	}
	currency := strings.TrimSpace(req.Currency)
	if currency == "" {
		return nil, gateway.ErrInvalidOrderCurrency
	}

	result, err := s.client.StartOrder(ctx, gateway.CreateOrderRequest{
		SagaID:               uuid.NewString(),
		OrderID:              uuid.NewString(),
		RequestedAt:          time.Now().UTC().Format(time.RFC3339Nano),
		IdempotencyKey:       uuid.NewString(),
		BuyerID:              buyerID,
		BuyerEmail:           buyerEmail,
		RealtimeConnectionID: connectionID,
		GigID:                gigID,
		GigTitle:             gigTitle,
		PackageID:            packageID,
		PackageTier:          packageTier,
		PackageDescription:   packageDescription,
		PackageDeliveryDays:  req.PackageDeliveryDays,
		PriceCents:           req.PriceCents,
		Currency:             currency,
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
		result.SagaID = req.SagaID
	}
	if strings.TrimSpace(result.OrderID) == "" {
		result.OrderID = req.OrderID
	}
	if strings.TrimSpace(result.Status) == "" {
		result.Status = "pending"
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

func validRealtimeConnectionID(id string) bool {
	parts := strings.Split(strings.TrimSpace(id), ".")
	return len(parts) == 2 && strings.TrimSpace(parts[0]) != "" && strings.TrimSpace(parts[1]) != ""
}
