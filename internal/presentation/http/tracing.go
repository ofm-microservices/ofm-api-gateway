package http

import (
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
		ctx := otel.GetTextMapPropagator().Extract(c.UserContext(), fiberCarrier{h: &c.Request().Header})
		ctx, span := tracer.Start(ctx, c.Route().Path,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", c.Method()),
				attribute.String("http.route", c.Route().Path),
			),
		)
		defer span.End()
		c.SetUserContext(ctx)

		err := c.Next()
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.SetAttributes(
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
