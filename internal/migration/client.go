package migration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	gateway "api-gateway/internal/domain"
	"github.com/google/uuid"
)

type recoveryMetadataContextKey struct{}

// RecoveryMetadata carries durable command identity from a Gateway request
// into a compatibility write without coupling the migration client to Fiber.
type RecoveryMetadata struct {
	CommandID      string
	CorrelationID  string
	IdempotencyKey string
}

// WithRecoveryMetadata attaches fallback identity to a request context.
func WithRecoveryMetadata(ctx context.Context, metadata RecoveryMetadata) context.Context {
	return context.WithValue(ctx, recoveryMetadataContextKey{}, metadata)
}

// LegacyRegistrationClient calls the monolith compatibility registration API.
type LegacyRegistrationClient interface {
	Start(ctx context.Context, req gateway.SignUpRequest) (*gateway.SignUpResult, error)
	Verify(ctx context.Context, req gateway.VerifyEmailRequest) (*gateway.VerifyEmailResult, error)
	Complete(ctx context.Context, req gateway.CompleteRegistrationRequest) (*gateway.AuthTokensResult, error)
}

type legacyClient struct {
	baseURL string
	client  *http.Client
}

// NewLegacyRegistrationClient constructs the monolith registration adapter.
func NewLegacyRegistrationClient(baseURL string, timeout time.Duration) (LegacyRegistrationClient, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return nil, fmt.Errorf("monolith base URL is empty")
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &legacyClient{baseURL: baseURL, client: &http.Client{Timeout: timeout}}, nil
}

func (c *legacyClient) Start(ctx context.Context, req gateway.SignUpRequest) (*gateway.SignUpResult, error) {
	var result gateway.SignUpResult
	if err := c.do(ctx, http.MethodPost, "/api/v2/auth/sign-up", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *legacyClient) Verify(ctx context.Context, req gateway.VerifyEmailRequest) (*gateway.VerifyEmailResult, error) {
	var result gateway.VerifyEmailResult
	if err := c.do(ctx, http.MethodPost, "/api/v2/auth/sign-up/verify-email", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *legacyClient) Complete(ctx context.Context, req gateway.CompleteRegistrationRequest) (*gateway.AuthTokensResult, error) {
	var result gateway.AuthTokensResult
	if err := c.do(ctx, http.MethodPost, "/api/v2/auth/sign-up/complete", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *legacyClient) do(ctx context.Context, method, path string, payload any, output any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal legacy request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create legacy request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	metadata, _ := ctx.Value(recoveryMetadataContextKey{}).(RecoveryMetadata)
	commandID := strings.TrimSpace(metadata.CommandID)
	if commandID == "" {
		commandID = uuid.NewString()
	}
	correlationID := strings.TrimSpace(metadata.CorrelationID)
	if correlationID == "" {
		correlationID = commandID
	}
	idempotencyKey := strings.TrimSpace(metadata.IdempotencyKey)
	if idempotencyKey == "" {
		idempotencyKey = commandID
	}
	req.Header.Set("X-Command-ID", commandID)
	req.Header.Set("X-Correlation-ID", correlationID)
	req.Header.Set("Idempotency-Key", idempotencyKey)
	response, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("legacy request failed: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		data, _ := io.ReadAll(io.LimitReader(response.Body, 64<<10))
		return fmt.Errorf("legacy request returned status %d: %s", response.StatusCode, strings.TrimSpace(string(data)))
	}
	if err := json.NewDecoder(response.Body).Decode(output); err != nil {
		return fmt.Errorf("decode legacy response: %w", err)
	}
	return nil
}
