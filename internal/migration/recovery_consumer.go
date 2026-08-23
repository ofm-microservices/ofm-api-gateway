package migration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	commonjwt "github.com/ofm-microservices/ofm-common/pkg/jwt"
	"github.com/ofm-microservices/ofm-common/pkg/migration/events"
	"github.com/segmentio/kafka-go"
)

func splitBrokers(raw string) []string {
	parts := strings.FieldsFunc(raw, func(r rune) bool { return r == ',' || r == ' ' || r == '\n' || r == '\t' })
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			result = append(result, part)
		}
	}
	return result
}

// RecoveryConsumer replays durable monolith fallback commands through the
// gateway's normal HTTP-to-gRPC routes. This keeps recovery adapters out of
// every bounded context while preserving the ordinary application flow.
type RecoveryConsumer interface {
	Run(context.Context) error
	Close() error
}

type recoveryConsumer struct {
	reader      *kafka.Reader
	dlq         *kafka.Writer
	completed   *kafka.Writer
	base        string
	client      *http.Client
	maxAttempts int
}

// NewRecoveryConsumer constructs a consumer for the monolith recovery topic.
func NewRecoveryConsumer(brokers, topic, group, dlq, baseURL string, timeout time.Duration, completedTopic ...string) (RecoveryConsumer, error) {
	if strings.TrimSpace(brokers) == "" || strings.TrimSpace(topic) == "" || strings.TrimSpace(group) == "" {
		return nil, fmt.Errorf("recovery Kafka configuration is incomplete")
	}
	if strings.TrimSpace(baseURL) == "" {
		return nil, fmt.Errorf("recovery gateway base URL is empty")
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	addresses := splitBrokers(brokers)
	// A cold-started service may need longer than one short retry burst to
	// become Ready. Keep the retry bounded, but allow a realistic recovery
	// window before moving the command to the DLQ.
	maxAttempts := 10
	if raw := strings.TrimSpace(os.Getenv("MIGRATION_RECOVERY_MAX_ATTEMPTS")); raw != "" {
		if parsed, parseErr := strconv.Atoi(raw); parseErr == nil && parsed > 0 {
			maxAttempts = parsed
		}
	}
	var completed *kafka.Writer
	if len(completedTopic) > 0 && strings.TrimSpace(completedTopic[0]) != "" {
		completed = &kafka.Writer{Addr: kafka.TCP(addresses...), Topic: completedTopic[0], Balancer: &kafka.Hash{}, RequiredAcks: kafka.RequireAll}
	}
	return &recoveryConsumer{
		reader:      kafka.NewReader(kafka.ReaderConfig{Brokers: addresses, Topic: topic, GroupID: group, MinBytes: 1, MaxBytes: 16 << 20}),
		dlq:         &kafka.Writer{Addr: kafka.TCP(addresses...), Topic: dlq, Balancer: &kafka.Hash{}, RequiredAcks: kafka.RequireAll},
		completed:   completed,
		base:        strings.TrimRight(baseURL, "/"),
		client:      &http.Client{Timeout: timeout},
		maxAttempts: maxAttempts,
	}, nil
}

// Close releases Kafka and HTTP resources owned by the consumer.
func (c *recoveryConsumer) Close() error {
	if c == nil {
		return nil
	}
	if err := c.reader.Close(); err != nil {
		_ = c.dlq.Close()
		if c.completed != nil {
			_ = c.completed.Close()
		}
		return err
	}
	if err := c.dlq.Close(); err != nil {
		return err
	}
	if c.completed != nil {
		return c.completed.Close()
	}
	return nil
}

// Run replays commands and commits only after the normal gateway route has
// accepted the write. Failed messages are retried with backoff and then moved
// to the configured DLQ.
func (c *recoveryConsumer) Run(ctx context.Context) error {
	for {
		message, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(5 * time.Second):
			}
			continue
		}
		var envelope events.Envelope
		if err = json.Unmarshal(message.Value, &envelope); err != nil {
			for {
				if dlqErr := c.publishDLQ(ctx, message, err); dlqErr == nil {
					break
				}
				if !waitRecovery(ctx, 5*time.Second) {
					return ctx.Err()
				}
			}
			for {
				if err = c.reader.CommitMessages(ctx, message); err == nil {
					break
				}
				if !waitRecovery(ctx, 5*time.Second) {
					return ctx.Err()
				}
			}
			continue
		}
		attempt := 0
		for {
			attempt++
			replayErr := c.replay(ctx, envelope)
			if replayErr == nil {
				break
			}
			if !replayErr.transient || attempt >= c.maxAttempts {
				for {
					if err = c.publishDLQ(ctx, message, replayErr); err == nil {
						break
					}
					if !waitRecovery(ctx, 5*time.Second) {
						return ctx.Err()
					}
				}
				break
			}
			delay := time.Duration(attempt) * time.Second
			if delay > 30*time.Second {
				delay = 30 * time.Second
			}
			if !waitRecovery(ctx, delay) {
				return ctx.Err()
			}
		}
		for {
			if err = c.reader.CommitMessages(ctx, message); err == nil {
				break
			}
			if !waitRecovery(ctx, 5*time.Second) {
				return ctx.Err()
			}
		}
	}
}

func waitRecovery(ctx context.Context, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

type recoveryReplayError struct {
	status    int
	transient bool
	message   string
}

func (e *recoveryReplayError) Error() string { return e.message }

func (c *recoveryConsumer) replay(ctx context.Context, envelope events.Envelope) *recoveryReplayError {
	if envelope.CommandMethod == "" || envelope.CommandPath == "" {
		return &recoveryReplayError{message: fmt.Sprintf("recovery command %s has no HTTP route", envelope.CommandID)}
	}
	req, err := http.NewRequestWithContext(ctx, envelope.CommandMethod, c.base+envelope.CommandPath, bytes.NewReader(envelope.Payload))
	if err != nil {
		return &recoveryReplayError{message: err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")
	// Keep the replay out of the fallback middleware without changing the
	// ordinary JWT/authentication path used by the HTTP handler.
	req.Header.Set("X-Recovery-Loop-Guard", "true")
	req.Header.Set("X-Command-ID", envelope.CommandID)
	req.Header.Set("X-Correlation-ID", envelope.CorrelationID)
	req.Header.Set("Idempotency-Key", envelope.IdempotencyKey)
	if envelope.RecoveryPrincipalID != "" {
		req.Header.Set("X-Recovery-Principal-ID", envelope.RecoveryPrincipalID)
		if envelope.RecoveryUsername != "" {
			req.Header.Set("X-Recovery-Username", envelope.RecoveryUsername)
		}
		if envelope.RecoveryEmail != "" {
			req.Header.Set("X-Recovery-Email", envelope.RecoveryEmail)
		}
		if req.Header.Get("Authorization") == "" {
			if token := recoveryToken(envelope.RecoveryPrincipalID, envelope.RecoveryUsername, envelope.RecoveryEmail); token != "" {
				req.Header.Set("Authorization", "Bearer "+token)
			}
		}
	}
	response, err := c.client.Do(req)
	if err != nil {
		return &recoveryReplayError{transient: true, message: err.Error()}
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 16<<20))
	if err != nil {
		return &recoveryReplayError{transient: true, message: err.Error()}
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return &recoveryReplayError{status: response.StatusCode, transient: response.StatusCode >= 500 || response.StatusCode == http.StatusTooManyRequests, message: fmt.Sprintf("recovery replay returned status %d", response.StatusCode)}
	}
	if err := c.publishCompleted(ctx, envelope, body); err != nil {
		return &recoveryReplayError{transient: true, message: err.Error()}
	}
	return nil
}

func recoveryToken(principal, username, email string) string {
	secret := strings.TrimSpace(os.Getenv("JWT_ACCESS_SECRET"))
	if secret == "" || strings.TrimSpace(principal) == "" {
		return ""
	}
	signer, err := commonjwt.NewSigner(commonjwt.Config{Secret: secret})
	if err != nil {
		return ""
	}
	if username == "" {
		username = "recovery"
	}
	if email == "" {
		email = "recovery@invalid.local"
	}
	token, err := signer.Sign(commonjwt.Claims{Subject: principal, Username: username, Email: email, IssuedAt: time.Now().UTC().Unix(), ExpiresAt: time.Now().UTC().Add(10 * time.Minute).Unix()})
	if err != nil {
		return ""
	}
	return token
}

func (c *recoveryConsumer) publishCompleted(ctx context.Context, command events.Envelope, body []byte) error {
	if c.completed == nil {
		return nil
	}
	aggregateID := responseAggregateID(command.AggregateType, body)
	event := events.Envelope{EventID: uuid.NewString(), CommandID: command.CommandID, CorrelationID: command.CorrelationID, CausationID: command.EventID, IdempotencyKey: command.IdempotencyKey, EventType: "migration.recovery.completed", Operation: command.Operation, SchemaVersion: 1, AggregateType: command.AggregateType, AggregateID: aggregateID, SourceService: "api-gateway-recovery", OccurredAt: time.Now().UTC(), Payload: body}
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return c.completed.WriteMessages(ctx, kafka.Message{Key: []byte(command.CommandID), Value: data})
}

func responseAggregateID(aggregateType string, body []byte) string {
	var value any
	if json.Unmarshal(body, &value) != nil {
		return ""
	}
	keys := map[string][]string{"gig": {"gig_id"}, "order": {"order_id"}, "user": {"user_id", "id"}, "review": {"review_id"}, "chat": {"message_id", "file_id"}, "file": {"file_id"}, "payment": {"payment_id", "payment_intent_id"}, "registration": {"session_id", "registration_id"}}
	var walk func(any) string
	walk = func(node any) string {
		if object, ok := node.(map[string]any); ok {
			for _, key := range keys[aggregateType] {
				if raw, ok := object[key].(string); ok {
					if _, err := uuid.Parse(raw); err == nil {
						return raw
					}
				}
			}
			for _, child := range object {
				if found := walk(child); found != "" {
					return found
				}
			}
		}
		if list, ok := node.([]any); ok {
			for _, child := range list {
				if found := walk(child); found != "" {
					return found
				}
			}
		}
		return ""
	}
	return walk(value)
}

func (c *recoveryConsumer) publishDLQ(ctx context.Context, message kafka.Message, cause error) error {
	value := append([]byte(nil), message.Value...)
	return c.dlq.WriteMessages(ctx, kafka.Message{Key: message.Key, Value: value, Headers: []kafka.Header{{Key: "recovery-error", Value: []byte(cause.Error())}, {Key: "original-topic", Value: []byte(message.Topic)}}})
}
