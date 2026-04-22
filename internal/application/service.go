package service

import (
	gateway "api-gateway/internal/domain"
	"context"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"
	"net/mail"
	"strings"
)

type registrationService struct {
	client RegistrationPublisher
	log    Logger
}

// New constructs the application service responsible for starting
// registrations through the saga boundary.
func New(client RegistrationPublisher, log Logger) (RegistrationService, error) {
	if client == nil {
		return nil, ErrNilRegistrationClient
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	return &registrationService{
		client: client,
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
