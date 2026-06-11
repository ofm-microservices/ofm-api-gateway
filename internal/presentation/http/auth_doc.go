package http

import gateway "api-gateway/internal/domain"

var (
	_ = gateway.SignUpRequest{}
	_ = gateway.SignUpResult{}
	_ = gateway.VerifyEmailRequest{}
	_ = gateway.VerifyEmailResult{}
	_ = gateway.CompleteRegistrationRequest{}
	_ = gateway.AuthTokensResult{}
	_ = gateway.SignInRequest{}
	_ = gateway.RefreshTokensRequest{}
	_ = gateway.User{}
)

// swaggerSignUpDoc documents the public signup endpoint exposed by api-gateway.
//
// @Summary Register account
// @Description Starts a registration session and triggers the first registration saga slice.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body gateway.SignUpRequest true "Signup payload"
// @Success 202 {object} gateway.SignUpResult
// @Failure 400 {object} SignUpBadRequestResponse
// @Failure 409 {object} SignUpConflictResponse
// @Failure 500 {object} SignUpInternalErrorResponse
// @Router /auth/sign-up [post]
func swaggerSignUpDoc() {}

// swaggerVerifyEmailDoc documents the verification-code submission endpoint.
//
// @Summary Verify registration email
// @Description Submits the emailed verification code to continue the registration saga.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body gateway.VerifyEmailRequest true "Verification payload"
// @Success 202 {object} gateway.VerifyEmailResult
// @Failure 400 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /auth/sign-up/verify-email [post]
func swaggerVerifyEmailDoc() {}

// swaggerCompleteRegistrationDoc documents the registration completion endpoint.
//
// @Summary Complete registration
// @Description Exchanges a completed registration session for auth-owned tokens.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body gateway.CompleteRegistrationRequest true "Completion payload"
// @Success 200 {object} gateway.AuthTokensResult
// @Failure 400 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /auth/sign-up/complete [post]
func swaggerCompleteRegistrationDoc() {}

// swaggerSignInDoc documents the credential sign-in endpoint.
//
// @Summary Sign in
// @Description Authenticates a user and returns a token pair.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body gateway.SignInRequest true "Sign-in payload"
// @Success 200 {object} gateway.AuthTokensResult
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /auth/sign-in [post]
func swaggerSignInDoc() {}

// swaggerRefreshDoc documents the refresh-token rotation endpoint.
//
// @Summary Refresh token pair
// @Description Rotates an auth refresh token and returns a new token pair.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body gateway.RefreshTokensRequest true "Refresh payload"
// @Success 200 {object} gateway.AuthTokensResult
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /auth/refresh [post]
func swaggerRefreshDoc() {}

// swaggerSignOutDoc documents the token revocation endpoint.
//
// @Summary Sign out
// @Description Revokes the presented refresh token and ends the session.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body gateway.SignOutRequest true "Sign-out payload"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /auth/sign-out [post]
func swaggerSignOutDoc() {}

// swaggerMeDoc documents the authenticated profile endpoint.
//
// @Summary Get current user
// @Description Returns the authenticated user's public preview profile.
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} gateway.User
// @Failure 401 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /auth/me [get]
func swaggerMeDoc() {}
