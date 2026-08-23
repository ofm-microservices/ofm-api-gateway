package migration

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// FallbackWriteClient replays a write request against the monolith when the
// owning gRPC service is unavailable. It is deliberately transport-level: the
// monolith receives the same public request contract and can persist a durable
// recovery command after its local write.
type FallbackWriteClient interface {
	Replay(ctx context.Context, method, path string, body []byte, headers http.Header) (int, []byte, http.Header, error)
}

type fallbackWriteClient struct {
	baseURL     string
	bypassToken string
	client      *http.Client
}

// NewFallbackWriteClient constructs the transport adapter used by the gateway
// write fallback. Reads are never sent through this client.
func NewFallbackWriteClient(baseURL, bypassToken string, timeout time.Duration) (FallbackWriteClient, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return nil, fmt.Errorf("monolith base URL is empty")
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &fallbackWriteClient{baseURL: baseURL, bypassToken: strings.TrimSpace(bypassToken), client: &http.Client{Timeout: timeout}}, nil
}

func (c *fallbackWriteClient) Replay(ctx context.Context, method, path string, body []byte, headers http.Header) (int, []byte, http.Header, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return 0, nil, nil, fmt.Errorf("create fallback request: %w", err)
	}
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	if c.bypassToken != "" {
		req.Header.Set("X-Monolith-Rate-Limit-Bypass", c.bypassToken)
	}
	response, err := c.client.Do(req)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("fallback request failed: %w", err)
	}
	defer response.Body.Close()
	data, err := io.ReadAll(io.LimitReader(response.Body, 16<<20))
	if err != nil {
		return 0, nil, nil, fmt.Errorf("read fallback response: %w", err)
	}
	return response.StatusCode, data, response.Header.Clone(), nil
}

func isWriteRequest(ctx *fiber.Ctx) bool {
	switch ctx.Method() {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}
