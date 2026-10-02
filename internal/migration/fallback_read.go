package migration

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// FallbackReadClient performs an idempotent GET against the monolith.
type FallbackReadClient interface {
	Read(context.Context, string, http.Header) (int, []byte, http.Header, error)
}

type fallbackReadClient struct {
	baseURL     string
	bypassToken string
	client      *http.Client
}

// ReadBreaker limits repeated calls to an unavailable primary dependency.
type ReadBreaker interface {
	Allow(time.Time) bool
	Success()
	Failure(time.Time)
}

type readBreaker struct {
	mu        sync.Mutex
	failures  int
	openUntil time.Time
}

// NewReadBreaker creates a small in-memory breaker for one gateway process.
func NewReadBreaker() ReadBreaker { return &readBreaker{} }
func (b *readBreaker) Allow(now time.Time) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return !now.Before(b.openUntil)
}
func (b *readBreaker) Success() {
	b.mu.Lock()
	b.failures = 0
	b.openUntil = time.Time{}
	b.mu.Unlock()
}
func (b *readBreaker) Failure(now time.Time) {
	b.mu.Lock()
	b.failures++
	if b.failures >= 3 {
		b.openUntil = now.Add(2 * time.Second)
	}
	b.mu.Unlock()
}

// NewFallbackReadClient constructs the read-only compatibility adapter.
func NewFallbackReadClient(baseURL, bypassToken string, timeout time.Duration) (FallbackReadClient, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return nil, fmt.Errorf("monolith base URL is empty")
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	transport := &http.Transport{MaxConnsPerHost: 128, MaxIdleConns: 256, MaxIdleConnsPerHost: 128, IdleConnTimeout: 90 * time.Second}
	return &fallbackReadClient{baseURL: baseURL, bypassToken: strings.TrimSpace(bypassToken), client: &http.Client{Timeout: timeout, Transport: transport}}, nil
}

func (c *fallbackReadClient) Read(ctx context.Context, path string, headers http.Header) (int, []byte, http.Header, error) {
	if !strings.HasPrefix(path, "/api/") {
		path = "/api" + path
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, bytes.NewReader(nil))
	if err != nil {
		return 0, nil, nil, err
	}
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	if c.bypassToken != "" {
		req.Header.Set("X-Monolith-Rate-Limit-Bypass", c.bypassToken)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("legacy read failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20))
	if err != nil {
		return resp.StatusCode, nil, resp.Header.Clone(), err
	}
	return resp.StatusCode, body, resp.Header.Clone(), nil
}
