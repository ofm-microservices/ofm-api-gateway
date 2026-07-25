package grpc

import (
	gateway "api-gateway/internal/domain"
	authv1 "github.com/ofm-microservices/ofm-common/proto/auth/v1"
	registrationv1 "github.com/ofm-microservices/ofm-common/proto/registration/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type registrationMapper struct{}

type authSessionMapper struct{}

func newRegistrationMapper() RegistrationMapper {
	return &registrationMapper{}
}

func newAuthSessionMapper() AuthSessionMapper {
	return &authSessionMapper{}
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

func (m *registrationMapper) ToVerifyEmailRequest(req VerifyEmailRequest) *registrationv1.VerifyEmailRequest {
	return &registrationv1.VerifyEmailRequest{
		SessionId: req.SessionID,
		ClientId:  req.ClientID,
		Code:      req.Code,
	}
}

func (m *registrationMapper) ToVerifyEmailResult(res *registrationv1.VerifyEmailResponse) *VerifyEmailResult {
	return &VerifyEmailResult{
		SessionID: res.GetSessionId(),
		ClientID:  res.GetClientId(),
		Status:    res.GetStatus(),
	}
}

func (m *registrationMapper) ToRegistrationStatus(res *registrationv1.GetRegistrationStatusResponse) *gateway.RegistrationStatus {
	return &gateway.RegistrationStatus{
		SessionID: res.GetSessionId(),
		ClientID:  res.GetClientId(),
		UserID:    res.GetUserId(),
		Status:    res.GetStatus(),
	}
}

func (m *registrationMapper) ToCompleteRegistrationResult(res *authv1.IssueRegistrationTokensResponse) *AuthTokensResult {
	return &AuthTokensResult{
		UserID:       res.GetUserId(),
		AccessToken:  res.GetAccessToken(),
		RefreshToken: res.GetRefreshToken(),
		TokenType:    res.GetTokenType(),
		ExpiresIn:    res.GetExpiresIn(),
	}
}

func (m *authSessionMapper) ToSignInRequest(req gateway.SignInRequest) *authv1.SignInRequest {
	return &authv1.SignInRequest{
		Identifier: req.Identifier,
		Password:   req.Password,
	}
}

func (m *authSessionMapper) ToRefreshRequest(req gateway.RefreshTokensRequest) *authv1.RefreshRequest {
	return &authv1.RefreshRequest{
		RefreshToken: req.RefreshToken,
	}
}

func (m *authSessionMapper) ToSignOutRequest(req gateway.SignOutRequest) *authv1.SignOutRequest {
	return &authv1.SignOutRequest{
		RefreshToken: req.RefreshToken,
	}
}

func (m *authSessionMapper) ToSignInResult(res *authv1.SignInResponse) *AuthTokensResult {
	return &AuthTokensResult{
		UserID:       res.GetUserId(),
		AccessToken:  res.GetAccessToken(),
		RefreshToken: res.GetRefreshToken(),
		TokenType:    res.GetTokenType(),
		ExpiresIn:    res.GetExpiresIn(),
	}
}

func (m *authSessionMapper) ToRefreshResult(res *authv1.RefreshResponse) *AuthTokensResult {
	return &AuthTokensResult{
		UserID:       res.GetUserId(),
		AccessToken:  res.GetAccessToken(),
		RefreshToken: res.GetRefreshToken(),
		TokenType:    res.GetTokenType(),
		ExpiresIn:    res.GetExpiresIn(),
	}
}

func (m *authSessionMapper) ToError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return gateway.ErrInvalidCredentials
	}

	switch st.Code() {
	case codes.Unauthenticated:
		return gateway.ErrInvalidCredentials
	case codes.InvalidArgument:
		switch st.Message() {
		case "invalid identifier":
			return gateway.ErrInvalidIdentifier
		case "invalid password":
			return gateway.ErrInvalidPassword
		case "invalid refresh token":
			return gateway.ErrInvalidRefreshToken
		default:
			return gateway.ErrInvalidCredentials
		}
	default:
		return gateway.ErrInvalidCredentials
	}
}

func (m *authSessionMapper) ToRefreshError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return gateway.ErrFailedToRefreshTokens
	}

	switch st.Code() {
	case codes.Unauthenticated:
		return gateway.ErrInvalidCredentials
	case codes.InvalidArgument:
		switch st.Message() {
		case "invalid refresh token":
			return gateway.ErrInvalidRefreshToken
		default:
			return gateway.ErrInvalidRefreshToken
		}
	default:
		return gateway.ErrFailedToRefreshTokens
	}
}

func (m *authSessionMapper) ToSignOutError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return gateway.ErrFailedToSignOut
	}

	switch st.Code() {
	case codes.Unauthenticated:
		return gateway.ErrInvalidCredentials
	case codes.InvalidArgument:
		switch st.Message() {
		case "invalid refresh token":
			return gateway.ErrInvalidRefreshToken
		default:
			return gateway.ErrInvalidRefreshToken
		}
	default:
		return gateway.ErrFailedToSignOut
	}
}

func (m *registrationMapper) ToStartRegistrationError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return gateway.ErrFailedToStartRegistration
	}

	switch st.Message() {
	case "invalid session id":
		return gateway.ErrInvalidSessionID
	case "invalid client id":
		return gateway.ErrInvalidClientID
	case "invalid verification code":
		return gateway.ErrInvalidVerificationCode
	case "invalid email":
		return gateway.ErrInvalidEmail
	case "invalid username":
		return gateway.ErrInvalidUsername
	case "invalid password":
		return gateway.ErrInvalidPassword
	case "registration is not ready for this operation":
		return gateway.ErrRegistrationNotCompleted
	case "registration session not found":
		return gateway.ErrRegistrationNotCompleted
	default:
		if st.Code() != codes.InvalidArgument {
			return gateway.ErrFailedToStartRegistration
		}
		return gateway.ErrFailedToStartRegistration
	}
}

func (m *registrationMapper) ToRegistrationStatusError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return gateway.ErrFailedToCompleteRegistration
	}

	switch st.Message() {
	case "invalid session id":
		return gateway.ErrInvalidSessionID
	case "invalid client id":
		return gateway.ErrInvalidClientID
	case "registration is not ready for this operation", "registration session not found":
		return gateway.ErrRegistrationNotCompleted
	default:
		return gateway.ErrFailedToCompleteRegistration
	}
}
