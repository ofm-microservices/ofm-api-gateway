package http

import gateway "api-gateway/internal/domain"

var _ = gateway.StartFreelancerOnboardingResult{}

// swaggerStartFreelancerOnboardingDoc documents the Stripe Connect onboarding start endpoint.
//
// @Summary Start freelancer onboarding
// @Description Starts Stripe Connect onboarding for the authenticated freelancer.
// @Tags onboarding
// @Produce json
// @Security BearerAuth
// @Success 200 {object} gateway.StartFreelancerOnboardingResult
// @Failure 401 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /freelancer/onboarding/start [post]
func swaggerStartFreelancerOnboardingDoc() {}
