package nats

import (
	"api-gateway/config"
	gateway "api-gateway/internal/domain"
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"time"
)

type orderPublisher struct {
	*publisher
}

// NewOrderPublisher constructs the order saga publisher used by api-gateway.
func NewOrderPublisher(cfg config.NATSConfig, log Logger) (OrderPublisher, error) {
	if cfg.URL == "" {
		return nil, ErrEmptyNATSURL
	}
	if cfg.OrderSagaStartSubject == "" {
		return nil, ErrEmptySubject
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	opts := []nats.Option{
		nats.Name("api-gateway"),
		nats.MaxReconnects(-1),
	}

	if cfg.User != "" {
		opts = append(opts, nats.UserInfo(cfg.User, cfg.Password))
	}

	nc, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, WrapConnectToNATSError(err)
	}

	return &orderPublisher{
		publisher: &publisher{
			nc:      nc,
			subject: cfg.OrderSagaStartSubject,
			log:     log.With(logging.String("module", "nats-publisher"), logging.String("subject", cfg.OrderSagaStartSubject)),
		},
	}, nil
}

// StartOrder publishes the order start payload to the configured saga subject.
func (p *orderPublisher) StartOrder(ctx context.Context, event gateway.CreateOrderRequest) (*gateway.CreateOrderResult, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, WrapMarshalEventError(err)
	}

	p.log.Info("publishing order event",
		logging.String("buyer_id", event.BuyerID),
		logging.String("gig_id", event.GigID),
	)

	if err := p.nc.Publish(p.subject, payload); err != nil {
		return nil, WrapPublishToNATSError(p.subject, err)
	}

	flushCtx := ctx
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		flushCtx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
	}

	if err := p.nc.FlushWithContext(flushCtx); err != nil {
		return nil, WrapFlushNATSError(err)
	}

	return &gateway.CreateOrderResult{Status: "accepted"}, nil
}
