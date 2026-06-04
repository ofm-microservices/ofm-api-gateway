package appfx

import (
	service "api-gateway/internal/application"
	"github.com/ofm-microservices/ofm-common/pkg/logging"

	"go.uber.org/fx"
)

// ServiceModule wires application services into the FX graph.
var ServiceModule = fx.Options(
	fx.Provide(ProvideRegistrationService),
	fx.Provide(ProvideAuthSessionService),
	fx.Provide(ProvideAuthMeService),
	fx.Provide(ProvideGigService),
	fx.Provide(ProvideUserProfileService),
	fx.Provide(ProvideOrderService),
	fx.Provide(ProvideOrderPreviewService),
	fx.Provide(ProvidePaymentOnboardingService),
	fx.Provide(ProvideReviewService),
	fx.Provide(ProvideSearchService),
)

// ProvideRegistrationService constructs the registration application service.
func ProvideRegistrationService(
	pub service.RegistrationPublisher,
	tokens service.TokenIssuer,
	lg logging.Logger,
) (service.RegistrationService, error) {
	return service.New(pub, tokens, lg)
}

// ProvideAuthSessionService constructs the auth session application service.
func ProvideAuthSessionService(
	client service.AuthSessionClient,
	lg logging.Logger,
) (service.AuthSessionService, error) {
	return service.NewAuthSession(client, lg)
}

// ProvideAuthMeService constructs the application service that resolves the
// current authenticated user's preview payload.
func ProvideAuthMeService(
	client service.UserClient,
	lg logging.Logger,
) (service.AuthMeService, error) {
	return service.NewAuthMe(client, lg)
}

// ProvideGigService constructs the gig application service.
func ProvideGigService(
	client service.GigPublisher,
	review service.ReviewClient,
	user service.UserClient,
	lg logging.Logger,
) (service.GigService, error) {
	return service.NewGig(client, review, user, lg)
}

// ProvideUserProfileService constructs the composite public user profile application service.
func ProvideUserProfileService(
	gigs service.GigService,
	reviews service.ReviewService,
	users service.UserClient,
	lg logging.Logger,
) (service.UserProfileService, error) {
	return service.NewUserProfile(gigs, reviews, users, lg)
}

// ProvideOrderService constructs the order application service.
func ProvideOrderService(
	client service.OrderCheckoutClient,
	lg logging.Logger,
) (service.OrderService, error) {
	return service.NewOrder(client, lg)
}

// ProvideOrderPreviewService constructs the user-scoped order preview service.
func ProvideOrderPreviewService(
	client service.OrderPreviewClient,
	payments service.PaymentByOrderClient,
	lg logging.Logger,
) (service.OrderPreviewService, error) {
	return service.NewOrderPreview(client, payments, lg)
}

// ProvidePaymentOnboardingService constructs the onboarding application service.
func ProvidePaymentOnboardingService(
	client service.PaymentOnboardingPublisher,
	lg logging.Logger,
) (service.PaymentOnboardingService, error) {
	return service.NewPaymentOnboarding(client, lg)
}

// ProvideReviewService constructs the review application service.
func ProvideReviewService(
	client service.ReviewClient,
	lg logging.Logger,
) (service.ReviewService, error) {
	return service.NewReview(client, lg)
}

// ProvideSearchService constructs the search application service.
func ProvideSearchService(
	client service.SearchClient,
	lg logging.Logger,
) (service.SearchService, error) {
	return service.NewSearch(client, lg)
}
