package migration

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	gateway "api-gateway/internal/domain"
)

func TestLegacyRegistrationClientStart(t *testing.T) {
	client, err := NewLegacyRegistrationClient("http://legacy.test", 0)
	if err != nil {
		t.Fatal(err)
	}
	client.(*legacyClient).client.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v2/auth/sign-up" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		return &http.Response{StatusCode: http.StatusAccepted, Body: io.NopCloser(strings.NewReader(`{"session_id":"session-1","client_id":"client-1","user_id":"user-1","status":"pending"}`)), Header: make(http.Header)}, nil
	})
	result, err := client.Start(context.Background(), gateway.SignUpRequest{Email: "user@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if result.SessionID != "session-1" || result.UserID != "user-1" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestLegacyRegistrationClientRejectsNonSuccess(t *testing.T) {
	client, err := NewLegacyRegistrationClient("http://legacy.test", 0)
	if err != nil {
		t.Fatal(err)
	}
	client.(*legacyClient).client.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusConflict, Body: io.NopCloser(strings.NewReader("bad request")), Header: make(http.Header)}, nil
	})
	if _, err := client.Start(context.Background(), gateway.SignUpRequest{}); err == nil {
		t.Fatal("expected legacy status error")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
