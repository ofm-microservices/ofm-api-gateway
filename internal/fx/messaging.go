package appfx

import (
	"api-gateway/config"
	service "api-gateway/internal/application"
	grpcclient "api-gateway/internal/presentation/grpc"
	"context"
	"github.com/ofm-microservices/ofm-common/pkg/logging"

	"go.uber.org/fx"
)

// MessagingModule provides the registration-saga client dependency.
var MessagingModule = fx.Options(
	fx.Provide(ProvideRegistrationPublisher),
	fx.Provide(ProvideTokenIssuer),
	fx.Provide(ProvideGigPublisher),
	fx.Provide(ProvideOrderPublisher),
	fx.Provide(ProvidePaymentOnboardingPublisher),
	fx.Provide(ProvideUserClient),
	fx.Provide(ProvideReviewClient),
	fx.Provide(ProvideSearchClient),
)

// ProvideRegistrationPublisher constructs the saga client and wires its
// lifecycle into FX shutdown.
func ProvideRegistrationPublisher(
	lc fx.Lifecycle,
	cfg *config.Config,
	lg logging.Logger,
) (service.RegistrationPublisher, error) {
	client, err := grpcclient.NewClient(cfg.RegistrationSaga, lg)
	if err != nil {
		lg.Error("connect registration saga grpc failed", logging.Err(err))
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			client.Close()
			return nil
		},
	})

	return client, nil
}

// ProvideTokenIssuer constructs the auth-service client used for final token
// exchange.
func ProvideTokenIssuer(
	lc fx.Lifecycle,
	cfg *config.Config,
	lg logging.Logger,
) (service.TokenIssuer, error) {
	client, err := grpcclient.NewAuthClient(cfg.AuthService, lg)
	if err != nil {
		lg.Error("connect auth service grpc failed", logging.Err(err))
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			client.Close()
			return nil
		},
	})

	return client, nil
}

// ProvideGigPublisher constructs the gig-service client used for draft and
// publish orchestration.
func ProvideGigPublisher(
	lc fx.Lifecycle,
	cfg *config.Config,
	lg logging.Logger,
) (service.GigPublisher, error) {
	client, err := grpcclient.NewGigClient(cfg.GigService, lg)
	if err != nil {
		lg.Error("connect gig service grpc failed", logging.Err(err))
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			client.Close()
			return nil
		},
	})

	return client, nil
}

// ProvideOrderPublisher constructs the order-checkout gRPC client used by the
// public order flow.
func ProvideOrderPublisher(
	lc fx.Lifecycle,
	cfg *config.Config,
	lg logging.Logger,
) (service.OrderCheckoutClient, error) {
	client, err := grpcclient.NewOrderCheckoutClient(cfg.OrderSaga, lg)
	if err != nil {
		lg.Error("connect order saga grpc failed", logging.Err(err))
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			client.Close()
			return nil
		},
	})

	return client, nil
}

// ProvidePaymentOnboardingPublisher constructs the payment-service gRPC client
// used for freelancer Stripe onboarding.
func ProvidePaymentOnboardingPublisher(
	lc fx.Lifecycle,
	cfg *config.Config,
	lg logging.Logger,
) (service.PaymentOnboardingPublisher, error) {
	client, err := grpcclient.NewPaymentOnboardingClient(cfg.PaymentService, lg)
	if err != nil {
		lg.Error("connect payment service grpc failed", logging.Err(err))
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			client.Close()
			return nil
		},
	})

	return client, nil
}

// ProvideUserClient constructs the user-service gRPC client used for public
// detailed freelancer profiles.
func ProvideUserClient(
	lc fx.Lifecycle,
	cfg *config.Config,
	lg logging.Logger,
) (service.UserClient, error) {
	client, err := grpcclient.NewUserClient(cfg.UserService, lg)
	if err != nil {
		lg.Error("connect user service grpc failed", logging.Err(err))
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			client.Close()
			return nil
		},
	})

	return client, nil
}

// ProvideReviewClient constructs the review-service gRPC client used for buyer reviews.
func ProvideReviewClient(
	lc fx.Lifecycle,
	cfg *config.Config,
	lg logging.Logger,
) (service.ReviewClient, error) {
	client, err := grpcclient.NewReviewClient(cfg.ReviewService, lg)
	if err != nil {
		lg.Error("connect review service grpc failed", logging.Err(err))
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			client.Close()
			return nil
		},
	})

	return client, nil
}

// ProvideSearchClient constructs the search-service gRPC client used for public search.
func ProvideSearchClient(
	lc fx.Lifecycle,
	cfg *config.Config,
	lg logging.Logger,
) (service.SearchClient, error) {
	client, err := grpcclient.NewSearchClient(cfg.SearchService, lg)
	if err != nil {
		lg.Error("connect search service grpc failed", logging.Err(err))
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			client.Close()
			return nil
		},
	})

	return client, nil
}
