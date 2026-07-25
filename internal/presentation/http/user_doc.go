package http

import gateway "api-gateway/internal/domain"

var _ = gateway.UserProfile{}

// swaggerGetUserByUsernameDoc documents the public user profile endpoint.
//
// @Summary Get user profile
// @Description Loads the public freelancer profile page by username.
// @Tags users
// @Produce json
// @Param username path string true "Username"
// @Param gigs_cursor query string false "Gig cursor"
// @Param reviews_cursor query string false "Review cursor"
// @Success 200 {object} gateway.UserProfile
// @Failure 400 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /users/{username} [get]
func swaggerGetUserByUsernameDoc() {}
