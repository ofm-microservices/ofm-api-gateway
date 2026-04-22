package grpc

import (
	gateway "api-gateway/internal/domain"
	registrationv1 "github.com/ofm-microseervices/ofm-common/proto/registration/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type registrationMapper struct{}

func newRegistrationMapper() RegistrationMapper {
	return &registrationMapper{}
}

func (m *registrationMapper) ToStartRegistrationRequest(req SignUpRequest) *registrationv1.StartRegistrationRequest {
	return &registrationv1.StartRegistrationRequest{
		ClientId:  req.ClientID,
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		Surname:   req.Surname,
	}
}

func (m *registrationMapper) ToSignUpResult(res *registrationv1.StartRegistrationResponse) *SignUpResult {
	return &SignUpResult{
		SessionID:     res.GetSessionId(),
		ClientID:      res.GetClientId(),
		UserID:        res.GetUserId(),
		Status:        res.GetStatus(),
		ConflictState: res.GetConflictState(),
		UsernameTaken: res.GetUsernameTaken(),
		EmailTaken:    res.GetEmailTaken(),
	}
}

func (m *registrationMapper) ToStartRegistrationError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return gateway.ErrFailedToStartRegistration
	}

	if st.Code() != codes.InvalidArgument {
		return gateway.ErrFailedToStartRegistration
	}

	switch st.Message() {
	case "invalid email":
		return gateway.ErrInvalidEmail
	case "invalid username":
		return gateway.ErrInvalidUsername
	case "invalid password":
		return gateway.ErrInvalidPassword
	default:
		return gateway.ErrFailedToStartRegistration
	}
}
