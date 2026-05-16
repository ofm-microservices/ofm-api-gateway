package nats

import (
	gateway "api-gateway/internal/domain"
	"context"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
)

// Logger aliases the shared logger contract used by the NATS adapter.
type Logger = logging.Logger

// RegistrationRequested aliases the signup payload published to NATS.
type RegistrationRequested = gateway.SignUpRequest

// Publisher starts registration by publishing the signup event.
type Publisher interface {
	StartRegistration(ctx context.Context, event RegistrationRequested) (*gateway.SignUpResult, error)
	Close()
}

// OrderRequest aliases the public order start payload published to NATS.
type OrderRequest = gateway.CreateOrderRequest

// OrderPublisher starts the order saga by publishing the order start command.
type OrderPublisher interface {
	StartOrder(ctx context.Context, event OrderRequest) (*gateway.CreateOrderResult, error)
	Close()
}
