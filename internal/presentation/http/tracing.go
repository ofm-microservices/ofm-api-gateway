package http

import (
	requestmetadata "github.com/ofm-microservices/ofm-common/pkg/observability/metadata"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type fiberCarrier struct {
	h *fasthttp.RequestHeader
}

func (c fiberCarrier) Get(key string) string { return string(c.h.Peek(key)) }
func (c fiberCarrier) Set(key, value string) { c.h.Set(key, value) }
func (c fiberCarrier) Keys() []string {
	keys := make([]string, 0)
	c.h.VisitAll(func(k, _ []byte) {
		keys = append(keys, string(k))
	})
	return keys
}

func tracingMiddleware() fiber.Handler {
	tracer := otel.Tracer("api-gateway/http")
	return func(c *fiber.Ctx) error {
		if c.Path() == "/metrics" {
			return c.Next()
		}

		ctx := otel.GetTextMapPropagator().Extract(c.UserContext(), fiberCarrier{h: &c.Request().Header})
		ctx, span := tracer.Start(ctx, c.Method()+" "+c.Path(),
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", c.Method()),
				attribute.String("http.target", c.Path()),
				attribute.String("ofm.test_run_id", c.Get("X-Test-Run-ID")),
				attribute.String("ofm.scenario_id", c.Get("X-Test-Scenario")),
				attribute.String("ofm.fault_profile", c.Get("X-Fault-Profile")),
				attribute.String("ofm.fault_target", c.Get("X-Fault-Target")),
			),
		)
		ctx = requestmetadata.WithValues(ctx, requestmetadata.Values{
			RequestID: c.Get("X-Request-ID"), CorrelationID: c.Get("X-Correlation-ID"), IdempotencyKey: c.Get("Idempotency-Key"), TestRunID: c.Get("X-Test-Run-ID"), TestScenarioID: c.Get("X-Test-Scenario"),
		})
		defer span.End()
		c.SetUserContext(ctx)

		err := c.Next()
		route := c.Route().Path
		if route == "" {
			route = "unknown"
		}
		span.SetName(route)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.SetAttributes(
			attribute.String("http.route", route),
			attribute.Int("http.status_code", c.Response().StatusCode()),
		)
		otel.GetTextMapPropagator().Inject(ctx, responseCarrier{h: &c.Response().Header})
		return err
	}
}

type responseCarrier struct {
	h *fasthttp.ResponseHeader
}

func (c responseCarrier) Get(key string) string { return string(c.h.Peek(key)) }
func (c responseCarrier) Set(key, value string) { c.h.Set(key, value) }
func (c responseCarrier) Keys() []string {
	keys := make([]string, 0)
	c.h.VisitAll(func(k, _ []byte) {
		keys = append(keys, string(k))
	})
	return keys
}
