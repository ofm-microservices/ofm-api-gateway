package service

import (
	gateway "api-gateway/internal/domain"
	"context"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"net/mail"
	"strings"
)

type registrationService struct {
	client RegistrationPublisher
	tokens TokenIssuer
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

func (s *registrationService) SignUp(ctx context.Context, req gateway.SignUpRequest) (*gateway.SignUpResult, error) {
	s.log.Info("sign up request received", logging.String("email", req.Email), logging.String("username", req.Username))

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
		s.log.Error("failed to start registration", logging.String("email", req.Email), logging.Err(err))
		return nil, err
	}
	if result.Status == "conflict" {
		return nil, &gateway.RegistrationConflictError{
			State:         result.ConflictState,
			UsernameTaken: result.UsernameTaken,
			EmailTaken:    result.EmailTaken,
		}
	}

	s.log.Info("registration started", logging.String("email", req.Email), logging.String("session_id", result.SessionID))
	return result, nil
}

func (s *registrationService) VerifyEmail(ctx context.Context, req gateway.VerifyEmailRequest) (*gateway.VerifyEmailResult, error) {
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
		s.log.Error("failed to verify registration email", logging.String("session_id", sessionID), logging.Err(err))
		return nil, err
	}

	return result, nil
}

func (s *registrationService) CompleteRegistration(ctx context.Context, req gateway.CompleteRegistrationRequest) (*gateway.CompleteRegistrationResult, error) {
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
		s.log.Error("failed to get registration status", logging.String("session_id", sessionID), logging.Err(err))
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
		s.log.Error("failed to issue registration tokens", logging.String("user_id", status.UserID), logging.Err(err))
		return nil, err
	}

	return result, nil
}
