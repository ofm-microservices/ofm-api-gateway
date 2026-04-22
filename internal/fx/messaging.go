package appfx

import (
	"api-gateway/config"
	service "api-gateway/internal/application"
	registrationgrpc "api-gateway/internal/presentation/grpc"
	"context"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"

	"go.uber.org/fx"
)

// MessagingModule provides the registration-saga client dependency.
var MessagingModule = fx.Options(
	fx.Provide(ProvideRegistrationPublisher),
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
