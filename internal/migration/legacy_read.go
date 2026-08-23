package migration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// LegacyReadClient reads normalized public responses from the monolith.
type LegacyReadClient interface {
	UserProfile(context.Context, string, string, string) (json.RawMessage, int, error)
}

type legacyReadClient struct {
	baseURL string
	client  *http.Client
}

// NewLegacyReadClient constructs the read-only legacy adapter.
func NewLegacyReadClient(baseURL string, timeout time.Duration) (LegacyReadClient, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return nil, fmt.Errorf("monolith base URL is empty")
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &legacyReadClient{baseURL: baseURL, client: &http.Client{Timeout: timeout}}, nil
}

func (c *legacyReadClient) UserProfile(ctx context.Context, username, gigsCursor, reviewsCursor string) (json.RawMessage, int, error) {
	path := "/api/v2/users/" + url.PathEscape(username) + "?gigs_cursor=" + url.QueryEscape(gigsCursor) + "&reviews_cursor=" + url.QueryEscape(reviewsCursor)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, 0, err
	}
	response, err := c.client.Do(request)
	if err != nil {
		return nil, 0, fmt.Errorf("legacy read failed: %w", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return nil, response.StatusCode, err
	}
	return json.RawMessage(body), response.StatusCode, nil
}
