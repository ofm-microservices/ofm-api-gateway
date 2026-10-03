package http

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type fallbackReplayStub struct {
	called bool
}

func (s *fallbackReplayStub) Replay(_ context.Context, method, path string, body []byte, headers http.Header) (int, []byte, http.Header, error) {
	s.called = true
	if method != http.MethodPost || path != "/api/v2/test" || len(body) != 2 || headers.Get("X-Command-ID") == "" {
		return http.StatusInternalServerError, nil, nil, fiber.ErrInternalServerError
	}
	return http.StatusAccepted, []byte(`{"recovered":true}`), make(http.Header), nil
}

func TestFallbackWriteMiddlewareFallsBackOnReturnedFiveHundredError(t *testing.T) {
	stub := &fallbackReplayStub{}
	app := fiber.New()
	app.Use(fallbackWriteMiddleware(stub))
	app.Post("/api/v2/test", func(*fiber.Ctx) error {
		return fiber.NewError(http.StatusServiceUnavailable, "owner unavailable")
	})

	req, err := http.NewRequest(http.MethodPost, "/api/v2/test", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected fallback status %d, got %d", http.StatusAccepted, resp.StatusCode)
	}
	if !stub.called {
		t.Fatal("expected fallback replay to be called")
	}
}
