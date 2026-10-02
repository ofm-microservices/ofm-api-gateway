package http

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"api-gateway/internal/migration"
	"github.com/gofiber/fiber/v2"
)

type readFallbackStub struct {
	calls int
	path  string
}

func (s *readFallbackStub) Read(_ context.Context, path string, _ http.Header) (int, []byte, http.Header, error) {
	s.calls++
	s.path = path
	return http.StatusOK, []byte(`{"source":"monolith"}`), make(http.Header), nil
}

var _ migration.FallbackReadClient = (*readFallbackStub)(nil)

func TestFallbackReadMiddlewareRetriesThenUsesMonolith(t *testing.T) {
	stub := &readFallbackStub{}
	app := fiber.New()
	app.Use(fallbackReadMiddleware(stub))
	app.Get("/api/v2/search", func(c *fiber.Ctx) error { return fiber.NewError(http.StatusServiceUnavailable, "search unavailable") })
	req, err := http.NewRequest(http.MethodGet, "/api/v2/search?query=test", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if stub.calls != 1 {
		t.Fatalf("expected one monolith call, got %d", stub.calls)
	}
	if stub.path != "/api/v2/search?query=test" {
		t.Fatalf("unexpected fallback path %q", stub.path)
	}
}
