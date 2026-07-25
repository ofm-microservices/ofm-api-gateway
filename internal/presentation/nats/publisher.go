package nats

import (
	"api-gateway/config"
	gateway "api-gateway/internal/domain"
	"context"
	"encoding/json"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"time"

	"github.com/nats-io/nats.go"
)

type publisher struct {
	nc      *nats.Conn
	subject string
	log     Logger
}

// NewPublisher constructs the legacy NATS-based registration publisher.
func NewPublisher(cfg config.NATSConfig, log Logger) (Publisher, error) {
	if cfg.URL == "" {
		return nil, ErrEmptyNATSURL
	}
	if cfg.SagaCreateAuthSubject == "" {
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

	return &publisher{
		nc:      nc,
		subject: cfg.SagaCreateAuthSubject,
		log:     log.With(logging.String("module", "nats-publisher"), logging.String("subject", cfg.SagaCreateAuthSubject)),
	}, nil
}

// StartRegistration publishes the signup payload to the configured saga
// subject and waits for the publisher flush to complete.
func (p *publisher) StartRegistration(ctx context.Context, event gateway.SignUpRequest) (*gateway.SignUpResult, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, WrapMarshalEventError(err)
	}

	p.log.Info("publishing registration event", logging.String("email", event.Email))

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

	return &gateway.SignUpResult{Status: "pending"}, nil
}

// Close closes the underlying NATS connection.
func (p *publisher) Close() {
	if p != nil && p.nc != nil {
		p.log.Info("closing nats publisher")
		p.nc.Close()
	}
}
