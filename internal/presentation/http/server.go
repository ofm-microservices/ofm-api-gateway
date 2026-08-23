package http

import (
	"api-gateway/config"
	"api-gateway/internal/migration"
	"context"
	"fmt"
	commonjwt "github.com/ofm-microservices/ofm-common/pkg/jwt"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"github.com/ofm-microservices/ofm-common/pkg/observability/metrics"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
)

type server struct {
	app *fiber.App
	cfg config.HTTPConfig
	log logging.Logger
}

// NewServer constructs the Fiber HTTP server used by api-gateway.
func NewServer(cfg config.HTTPConfig, log logging.Logger, fallback ...migration.FallbackWriteClient) (Server, error) {
	if log == nil {
		return nil, ErrNilLogger
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://127.0.0.1:3000,http://api.ofm.local",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,Idempotency-Key",
		AllowCredentials: true,
	}))
	app.Use(tracingMiddleware())
	app.Use(metricsMiddleware)
	if len(fallback) > 0 && fallback[0] != nil {
		app.Use(fallbackWriteMiddleware(fallback[0]))
	}
	app.Use(newFaultInjector(cfg.FaultInjection, cfg.Environment, time.Now()).Middleware())
	srv := &server{
		app: app,
		cfg: cfg,
		log: log.With(logging.String("module", "http-server")),
	}

	return srv, nil
}

// fallbackWriteMiddleware retries failed v2 writes through the monolith. It
// only handles dependency failures after the normal gRPC route has run; reads,
// validation errors and successful microservice responses remain untouched.
func fallbackWriteMiddleware(client migration.FallbackWriteClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		loopGuard := c.Get("X-Recovery-Replay") == "true" || c.Get("X-Recovery-Loop-Guard") == "true"
		if !isFallbackWritePath(c.Path()) || !isWriteMethod(c.Method()) {
			return c.Next()
		}
		requestBody := append([]byte(nil), c.Body()...)
		err := c.Next()
		statusCode := c.Response().StatusCode()
		if err != nil {
			if fiberErr, ok := err.(*fiber.Error); ok && fiberErr.Code >= 100 {
				statusCode = fiberErr.Code
			}
		}
		if loopGuard || statusCode < 500 {
			return err
		}
		headers := make(http.Header)
		commandID := c.Get("X-Command-ID")
		if commandID == "" {
			commandID = uuid.NewString()
		}
		correlationID := c.Get("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.NewString()
		}
		for _, key := range []string{"Authorization", "Content-Type", "Idempotency-Key", "X-Request-ID", "X-Correlation-ID", "Traceparent"} {
			if value := c.Get(key); value != "" {
				headers.Set(key, value)
			}
		}
		headers.Set("X-Command-ID", commandID)
		headers.Set("X-Correlation-ID", correlationID)
		// The replay is an internal recovery request. The gateway must use the
		// explicit recovery principal instead of trying to parse a bearer token
		// that is intentionally not transported through the recovery queue.
		headers.Set("X-Recovery-Replay", "true")
		if principal, ok := c.Locals("jwt.freelancer_id").(string); ok && principal != "" {
			headers.Set("X-Recovery-Principal-ID", principal)
		}
		if claims, ok := c.Locals(jwtClaimsLocalKey).(*commonjwt.Claims); ok && claims != nil {
			if claims.Username != "" {
				headers.Set("X-Recovery-Username", claims.Username)
			}
			if claims.Email != "" {
				headers.Set("X-Recovery-Email", claims.Email)
			}
		}
		if headers.Get("Idempotency-Key") == "" {
			headers.Set("Idempotency-Key", commandID)
		}
		status, body, responseHeaders, fallbackErr := client.Replay(c.UserContext(), c.Method(), c.OriginalURL(), requestBody, headers)
		if fallbackErr != nil {
			return err
		}
		for key, values := range responseHeaders {
			if key == "Content-Length" || key == "Transfer-Encoding" {
				continue
			}
			for _, value := range values {
				c.Set(key, value)
			}
		}
		c.Status(status)
		return c.Send(body)
	}
}

func isFallbackWritePath(path string) bool {
	return strings.HasPrefix(path, "/api/v2/")
}

func isWriteMethod(method string) bool {
	return method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete
}

func metricsMiddleware(c *fiber.Ctx) error {
	meter := metrics.Global()
	meter.IncHTTPInFlight()
	started := time.Now()
	err := c.Next()
	route := c.Route().Path
	if route == "" {
		route = "unknown"
	}
	meter.ObserveHTTP(c.Method(), route, metrics.HTTPStatusClass(c.Response().StatusCode()), time.Since(started), len(c.Body()), len(c.Response().Body()))
	meter.DecHTTPInFlight()
	return err
}

// App returns the underlying Fiber application for route registration.
func (s *server) App() *fiber.App {
	return s.app
}

// Start begins serving the public HTTP API.
func (s *server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	s.log.Info("starting http server", logging.String("addr", addr))
	return s.app.Listen(addr)
}

// Shutdown stops the HTTP server using the provided context deadline.
func (s *server) Shutdown(ctx context.Context) error {
	s.log.Info("shutting down http server")
	return s.app.ShutdownWithContext(ctx)
}
