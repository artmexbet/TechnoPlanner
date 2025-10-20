package app

import (
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// TelemetryMiddleware creates a middleware for OpenTelemetry tracing
func TelemetryMiddleware(serviceName string) fiber.Handler {
	tracer := otel.Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(c *fiber.Ctx) error {
		// Extract context from incoming request headers
		ctx := propagator.Extract(c.UserContext(), &fiberCarrier{ctx: c})

		// Start a new span
		ctx, span := tracer.Start(ctx, c.Path(),
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPMethodKey.String(c.Method()),
				semconv.HTTPRouteKey.String(c.Route().Path),
				semconv.HTTPTargetKey.String(string(c.Request().RequestURI())),
				semconv.HTTPURLKey.String(c.OriginalURL()),
				attribute.String("http.client_ip", c.IP()),
			),
		)
		defer span.End()

		// Store context in fiber context
		c.SetUserContext(ctx)

		// Continue processing
		err := c.Next()

		// Record response status
		span.SetAttributes(
			semconv.HTTPStatusCodeKey.Int(c.Response().StatusCode()),
		)

		if err != nil {
			span.RecordError(err)
		}

		return err
	}
}

// fiberCarrier adapts fiber.Ctx to satisfy the TextMapCarrier interface
type fiberCarrier struct {
	ctx *fiber.Ctx
}

func (fc *fiberCarrier) Get(key string) string {
	return fc.ctx.Get(key)
}

func (fc *fiberCarrier) Set(key string, value string) {
	fc.ctx.Set(key, value)
}

func (fc *fiberCarrier) Keys() []string {
	keys := make([]string, 0)
	fc.ctx.Request().Header.VisitAll(func(key, _ []byte) {
		keys = append(keys, string(key))
	})
	return keys
}
