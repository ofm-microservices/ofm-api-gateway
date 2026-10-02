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
func NewServer(cfg config.HTTPConfig, log logging.Logger, dependencies ...any) (Server, error) {
	if log == nil {
		return nil, ErrNilLogger
	}

	app := fiber.New()
	app.Use(transportLoggingMiddleware(log))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://127.0.0.1:3000,http://api.ofm.local",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,Idempotency-Key",
		AllowCredentials: true,
	}))
	app.Use(tracingMiddleware())
	app.Use(metricsMiddleware)
	for _, dependency := range dependencies {
		switch client := dependency.(type) {
		case migration.FallbackWriteClient:
			app.Use(fallbackWriteMiddlewareWithLogger(client, log, cfg.JWTAccessSecret))
		case migration.FallbackReadClient:
			app.Use(fallbackReadMiddlewareWithLogger(client, log))
		}
	}
	app.Use(newFaultInjector(cfg.FaultInjection, cfg.Environment, time.Now()).Middleware())
	srv := &server{
		app: app,
		cfg: cfg,
		log: log.With(logging.String("module", "http-server")),
	}

	return srv, nil
}

// fallbackReadMiddleware retries idempotent dependency failures and then uses
// the monolith compatibility API. It never creates recovery commands.
func fallbackReadMiddleware(client migration.FallbackReadClient) fiber.Handler {
	return fallbackReadMiddlewareWithLogger(client, nil)
}

func fallbackReadMiddlewareWithLogger(client migration.FallbackReadClient, log logging.Logger) fiber.Handler {
	breaker := migration.NewReadBreaker()
	return func(c *fiber.Ctx) error {
		if c.Method() != http.MethodGet || !isReadFallbackPath(c.Path()) || c.Get("X-Read-Fallback-Attempt") != "" {
			return c.Next()
		}
		var err error
		for attempt := 0; attempt < 2 && breaker.Allow(time.Now()); attempt++ {
			err = c.Next()
			if !readDependencyFailure(c.Response().StatusCode(), err) {
				breaker.Success()
				return err
			}
			breaker.Failure(time.Now())
			c.Response().Reset()
			time.Sleep(time.Duration(attempt+1) * 25 * time.Millisecond)
		}
		headers := make(http.Header)
		for _, key := range []string{"Authorization", "X-Request-ID", "X-Correlation-ID", "Traceparent", "X-Test-Run-ID", "X-Test-Scenario"} {
			if value := c.Get(key); value != "" {
				headers.Set(key, value)
			}
		}
		if claims, ok := c.Locals(jwtClaimsLocalKey).(*commonjwt.Claims); ok && claims != nil {
			principal := strings.TrimSpace(claims.Subject)
			if localPrincipal, localOK := c.Locals(jwtPrincipalLocalKey).(string); localOK && strings.TrimSpace(localPrincipal) != "" {
				principal = strings.TrimSpace(localPrincipal)
			}
			if principal != "" {
				headers.Set("X-Recovery-Principal-ID", principal)
			}
			if claims.Username != "" {
				headers.Set("X-Recovery-Username", claims.Username)
			}
			if claims.Email != "" {
				headers.Set("X-Recovery-Email", claims.Email)
			}
		}
		headers.Set("X-Read-Fallback-Attempt", "true")
		requestPath := c.Path()
		if query := string(c.Request().URI().QueryString()); query != "" {
			requestPath += "?" + query
		}
		status, body, responseHeaders, fallbackErr := client.Read(context.WithoutCancel(c.UserContext()), requestPath, headers)
		if log != nil {
			fields := []logging.Field{
				logging.String("path", requestPath),
				logging.Int("primary_status", c.Response().StatusCode()),
				logging.Int("fallback_status", status),
				logging.String("run_id", c.Get("X-Test-Run-ID")),
				logging.String("scenario_id", c.Get("X-Test-Scenario")),
			}
			if fallbackErr != nil {
				fields = append(fields, logging.Err(fallbackErr))
				log.Error("read fallback failed", append(fields, logging.Operation("migration.read_fallback.failed"))...)
			} else {
				log.Info("read fallback completed", append(fields, logging.Operation("migration.read_fallback.completed"))...)
			}
		}
		if fallbackErr != nil || status >= 500 {
			if err == nil {
				err = fiber.ErrServiceUnavailable
			}
			return err
		}
		for key, values := range responseHeaders {
			if key != "Content-Length" && key != "Transfer-Encoding" {
				for _, value := range values {
					c.Set(key, value)
				}
			}
		}
		c.Status(status).Set("X-Read-Fallback", "monolith")
		return c.Send(body)
	}
}

func isReadFallbackPath(path string) bool {
	return strings.HasPrefix(path, "/v1/") || strings.HasPrefix(path, "/v2/") || strings.HasPrefix(path, "/api/v1/") || strings.HasPrefix(path, "/api/v2/")
}

func readDependencyFailure(status int, err error) bool {
	if err != nil {
		return true
	}
	return status == http.StatusInternalServerError || status == http.StatusBadGateway || status == http.StatusServiceUnavailable || status == http.StatusGatewayTimeout
}

// fallbackWriteMiddleware retries failed v2 writes through the monolith. It
// only handles dependency failures after the normal gRPC route has run; reads,
// validation errors and successful microservice responses remain untouched.
func fallbackWriteMiddleware(client migration.FallbackWriteClient, jwtSecret ...string) fiber.Handler {
	return fallbackWriteMiddlewareWithLogger(client, nil, jwtSecret...)
}

func fallbackWriteMiddlewareWithLogger(client migration.FallbackWriteClient, log logging.Logger, jwtSecret ...string) fiber.Handler {
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
		// Preserve the authenticated bearer for the internal replay. The local
		// migration monolith uses the same JWT contract; recovery headers remain
		// available as explicit identity metadata and the loop guard prevents
		// replaying the request back through fallback again.
		for _, key := range []string{"Authorization", "Content-Type", "Idempotency-Key", "X-Request-ID", "X-Correlation-ID", "X-Test-Run-ID", "X-Test-Scenario", "Traceparent"} {
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
		principal := ""
		if len(jwtSecret) > 0 {
			principal = principalFromAuthorization(c.Get("Authorization"), jwtSecret[0])
		}
		if principal != "" {
			headers.Set("X-Recovery-Principal-ID", principal)
		}
		if claims, ok := c.Locals(jwtClaimsLocalKey).(*commonjwt.Claims); ok && claims != nil {
			principal = strings.TrimSpace(claims.Subject)
			if localPrincipal, localOK := c.Locals(jwtPrincipalLocalKey).(string); localOK && strings.TrimSpace(localPrincipal) != "" {
				principal = strings.TrimSpace(localPrincipal)
			}
			if principal != "" {
				headers.Set("X-Recovery-Principal-ID", principal)
			}
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
		// The primary dependency call may have canceled the request context after
		// its timeout. Recovery is a new bounded attempt and must not inherit
		// that cancellation; WithoutCancel preserves request values used for
		// tracing and correlation while removing the expired deadline.
		recoveryCtx := context.WithoutCancel(c.UserContext())
		status, body, responseHeaders, fallbackErr := client.Replay(recoveryCtx, c.Method(), c.OriginalURL(), requestBody, headers)
		if fallbackErr != nil {
			if log != nil {
				log.Error("legacy fallback replay failed", logging.Int("fallback_status", status), logging.String("fallback_body", string(body)), logging.String("path", c.OriginalURL()), logging.String("run_id", c.Get("X-Test-Run-ID")), logging.Err(fallbackErr))
			}
			return err
		}
		if status < http.StatusOK || status >= http.StatusMultipleChoices {
			if log != nil {
				log.Error("legacy fallback replay returned non-success", logging.Int("fallback_status", status), logging.String("fallback_body", string(body)), logging.String("path", c.OriginalURL()), logging.String("run_id", c.Get("X-Test-Run-ID")))
			}
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
		// A fallback is accepted only when the monolith completed the command.
		// Validation/auth/domain failures must remain failures and must not be
		// counted as successful migration fallback evidence.
		if status >= http.StatusOK && status < http.StatusBadRequest {
			c.Set("X-Fallback-Accepted", "true")
		}
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
