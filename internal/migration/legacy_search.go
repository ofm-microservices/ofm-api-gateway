package migration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	gateway "api-gateway/internal/domain"
)

// LegacySearchClient is the read-only HTTP adapter used while search is
// migrated from the search service to the monolith PostgreSQL projection.
type LegacySearchClient interface {
	Search(context.Context, gateway.SearchRequest) (*gateway.SearchResponse, error)
	Close() error
}

type legacySearchClient struct {
	baseURL string
	client  *http.Client
}

// NewLegacySearchClient constructs the monolith-backed search adapter.
func NewLegacySearchClient(baseURL string, timeout time.Duration) (LegacySearchClient, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return nil, fmt.Errorf("monolith base URL is empty")
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &legacySearchClient{baseURL: baseURL, client: &http.Client{Timeout: timeout}}, nil
}

func (c *legacySearchClient) Search(ctx context.Context, req gateway.SearchRequest) (*gateway.SearchResponse, error) {
	query := url.Values{}
	query.Set("query", req.Query)
	query.Set("sort", strconv.FormatInt(int64(req.Sort), 10))
	query.Set("order", strconv.FormatInt(int64(req.Order), 10))
	if req.Cursor != "" {
		query.Set("cursor", req.Cursor)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v2/search?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}
	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("legacy search failed: %w", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("legacy search returned status %d: %s", response.StatusCode, strings.TrimSpace(string(body)))
	}
	var result gateway.SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode legacy search response: %w", err)
	}
	return &result, nil
}

func (c *legacySearchClient) Close() error { return nil }
