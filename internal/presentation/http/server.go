package http

import (
	"api-gateway/config"
	"context"
	"fmt"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"github.com/ofm-microservices/ofm-common/pkg/observability/metrics"
	"time"

	"github.com/gofiber/fiber/v2"
)

type server struct {
	app *fiber.App
	cfg config.HTTPConfig
	log logging.Logger
}

// NewServer constructs the Fiber HTTP server used by api-gateway.
func NewServer(cfg config.HTTPConfig, log logging.Logger) (Server, error) {
	if log == nil {
		return nil, ErrNilLogger
	}

	app := fiber.New()
	app.Use(tracingMiddleware())
	app.Use(metricsMiddleware)
	srv := &server{
		app: app,
		cfg: cfg,
		log: log.With(logging.String("module", "http-server")),
	}

	return srv, nil
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
