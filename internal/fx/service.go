package appfx

import (
	service "api-gateway/internal/application"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"

	"go.uber.org/fx"
)

// ServiceModule wires application services into the FX graph.
var ServiceModule = fx.Options(
	fx.Provide(ProvideRegistrationService),
)

// ProvideRegistrationService constructs the registration application service.
func ProvideRegistrationService(pub service.RegistrationPublisher, lg logging.Logger) (service.RegistrationService, error) {
	return service.New(pub, lg)
}
