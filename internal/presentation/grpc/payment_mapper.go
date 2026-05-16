package grpc

import (
	gateway "api-gateway/internal/domain"
	paymentconnectv1 "github.com/ofm-microservices/ofm-common/proto/paymentconnect/v1"
)

type paymentMapper struct{}

func newPaymentMapper() *paymentMapper { return &paymentMapper{} }

func (m *paymentMapper) ToStartFreelancerOnboardingRequest(req gateway.StartFreelancerOnboardingRequest) *paymentconnectv1.StartFreelancerOnboardingRequest {
	return &paymentconnectv1.StartFreelancerOnboardingRequest{
		UserId:         req.UserID,
		Country:        req.Country,
		ReturnUrl:      req.ReturnURL,
		RefreshUrl:     req.RefreshURL,
		IdempotencyKey: "",
		RequestedAt:    "",
	}
}

func (m *paymentMapper) ToStartFreelancerOnboardingResponse(res *paymentconnectv1.StartFreelancerOnboardingResponse) *gateway.StartFreelancerOnboardingResult {
	return &gateway.StartFreelancerOnboardingResult{
		UserID:           res.GetUserId(),
		StripeAccountID:  res.GetStripeAccountId(),
		OnboardingURL:    res.GetOnboardingUrl(),
		Status:           res.GetStatus(),
		DetailsSubmitted: res.GetDetailsSubmitted(),
		ChargesEnabled:   res.GetChargesEnabled(),
		PayoutsEnabled:   res.GetPayoutsEnabled(),
		DisabledReason:   res.GetDisabledReason(),
	}
}

func (m *paymentMapper) ToError(err error) error { return err }
