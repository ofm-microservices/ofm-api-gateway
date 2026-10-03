package http

import (
	"api-gateway/config"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestFaultInjectorRequiresTokenAndReturnsConfiguredStatus(t *testing.T) {
	app := fiber.New()
	injector := newFaultInjector(config.FaultInjectionConfig{
		Enabled:   true,
		Profile:   faultProfile503,
		Target:    "api-gateway",
		TestToken: "secret",
	}, "k3d", time.Now())
	app.Use(injector.Middleware())
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	withoutToken, err := app.Test(httptest.NewRequest("GET", "/health", nil))
	if err != nil || withoutToken.StatusCode != fiber.StatusOK {
		t.Fatalf("request without token: status=%d err=%v", withoutToken.StatusCode, err)
	}

	request := httptest.NewRequest("GET", "/health", nil)
	request.Header.Set("X-Fault-Test-Token", "secret")
	request.Header.Set("X-Fault-Target", "api-gateway")
	withToken, err := app.Test(request)
	if err != nil || withToken.StatusCode != fiber.StatusServiceUnavailable {
		t.Fatalf("request with token: status=%d err=%v", withToken.StatusCode, err)
	}
}

func TestFaultInjectorIsDisabledInProduction(t *testing.T) {
	app := fiber.New()
	injector := newFaultInjector(config.FaultInjectionConfig{
		Enabled:   true,
		Profile:   faultProfile503,
		TestToken: "secret",
	}, "production", time.Now())
	app.Use(injector.Middleware())
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	request := httptest.NewRequest("GET", "/health", nil)
	request.Header.Set("X-Fault-Test-Token", "secret")
	response, err := app.Test(request)
	if err != nil || response.StatusCode != fiber.StatusOK {
		t.Fatalf("production request: status=%d err=%v", response.StatusCode, err)
	}
}
