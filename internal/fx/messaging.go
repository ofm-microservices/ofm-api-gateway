package appfx

import (
	"api-gateway/config"
	service "api-gateway/internal/application"
	paymentgrpc "api-gateway/internal/presentation/grpc"
	registrationgrpc "api-gateway/internal/presentation/grpc"
	natspub "api-gateway/internal/presentation/nats"
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
)

// ProvideRegistrationPublisher constructs the saga client and wires its
// lifecycle into FX shutdown.
func ProvideRegistrationPublisher(
	lc fx.Lifecycle,
	cfg *config.Config,
	lg logging.Logger,
) (service.RegistrationPublisher, error) {
	client, err := registrationgrpc.NewClient(cfg.RegistrationSaga, lg)
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
	client, err := registrationgrpc.NewAuthClient(cfg.AuthService, lg)
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
	client, err := registrationgrpc.NewGigClient(cfg.GigService, lg)
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

// ProvideOrderPublisher constructs the order-saga publisher used by the public
// order create endpoint.
func ProvideOrderPublisher(
	lc fx.Lifecycle,
	cfg *config.Config,
	lg logging.Logger,
) (service.OrderPublisher, error) {
	client, err := natspub.NewOrderPublisher(cfg.NATS, lg)
	if err != nil {
		lg.Error("connect order saga nats failed", logging.Err(err))
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
	client, err := paymentgrpc.NewPaymentOnboardingClient(cfg.PaymentService, lg)
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
