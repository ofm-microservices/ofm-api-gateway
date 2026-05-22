package appfx

import (
	service "api-gateway/internal/application"
	"github.com/ofm-microservices/ofm-common/pkg/logging"

	"go.uber.org/fx"
)

// ServiceModule wires application services into the FX graph.
var ServiceModule = fx.Options(
	fx.Provide(ProvideRegistrationService),
	fx.Provide(ProvideGigService),
	fx.Provide(ProvideOrderService),
	fx.Provide(ProvidePaymentOnboardingService),
)

// ProvideRegistrationService constructs the registration application service.
func ProvideRegistrationService(
	pub service.RegistrationPublisher,
	tokens service.TokenIssuer,
	lg logging.Logger,
) (service.RegistrationService, error) {
	return service.New(pub, tokens, lg)
}

// ProvideGigService constructs the gig application service.
func ProvideGigService(
	client service.GigPublisher,
	lg logging.Logger,
) (service.GigService, error) {
	return service.NewGig(client, lg)
}

// ProvideOrderService constructs the order application service.
func ProvideOrderService(
	client service.OrderCheckoutClient,
	lg logging.Logger,
) (service.OrderService, error) {
	return service.NewOrder(client, lg)
}

// ProvidePaymentOnboardingService constructs the onboarding application service.
func ProvidePaymentOnboardingService(
	client service.PaymentOnboardingPublisher,
	lg logging.Logger,
) (service.PaymentOnboardingService, error) {
	return service.NewPaymentOnboarding(client, lg)
}
