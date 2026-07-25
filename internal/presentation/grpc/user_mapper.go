package grpc

import (
	gateway "api-gateway/internal/domain"
	userv1 "github.com/ofm-microservices/ofm-common/proto/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userMapper struct{}

func newUserMapper() UserMapper { return &userMapper{} }

func (m *userMapper) ToGetUserPreviewRequest(userID string) *userv1.GetUserPreviewByIDRequest {
	return &userv1.GetUserPreviewByIDRequest{UserId: userID}
}

func (m *userMapper) ToGetUserPreviewResponse(res *userv1.GetUserPreviewByIDResponse) *gateway.User {
	if res == nil || res.GetUser() == nil {
		return nil
	}

	user := res.GetUser()
	return &gateway.User{
		UserID:      user.GetUserId(),
		Username:    user.GetUsername(),
		DisplayName: user.GetDisplayName(),
		AvatarID:    user.GetAvatarId(),
		AvatarURL:   user.GetAvatarUrl(),
	}
}

func (m *userMapper) ToGetDetailedUserRequest(username string) *userv1.GetDetailedUserByUsernameRequest {
	return &userv1.GetDetailedUserByUsernameRequest{Username: username}
}

func (m *userMapper) ToGetDetailedUserResponse(res *userv1.GetDetailedUserByUsernameResponse) *gateway.User {
	if res == nil || res.GetUser() == nil {
		return nil
	}

	user := res.GetUser()
	return &gateway.User{
		UserID:      user.GetUserId(),
		Username:    user.GetUsername(),
		DisplayName: user.GetDisplayName(),
		AvatarID:    user.GetAvatarId(),
		AvatarURL:   user.GetAvatarUrl(),
		About:       user.GetAbout(),
	}
}

func (m *userMapper) ToError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	switch st.Code() {
	case codes.NotFound:
		return gateway.ErrUserNotFound
	case codes.InvalidArgument:
		return gateway.ErrInvalidUsername
	default:
		return err
	}
}
