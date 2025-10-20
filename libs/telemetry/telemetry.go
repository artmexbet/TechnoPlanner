package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// Config holds the configuration for OpenTelemetry
type Config struct {
	ServiceName     string `yaml:"service_name" env:"SERVICE_NAME"`
	JaegerEndpoint  string `yaml:"jaeger_endpoint" env:"JAEGER_ENDPOINT"`
	Enabled         bool   `yaml:"enabled" env:"ENABLED"`
}

// Telemetry holds the OpenTelemetry tracer provider
type Telemetry struct {
	TracerProvider *sdktrace.TracerProvider
	Tracer         trace.Tracer
}

// New creates a new Telemetry instance with Jaeger exporter
func New(cfg Config) (*Telemetry, error) {
	if !cfg.Enabled {
		// Return a no-op telemetry if disabled
		return &Telemetry{
			TracerProvider: sdktrace.NewTracerProvider(),
			Tracer:         otel.Tracer(cfg.ServiceName),
		}, nil
	}

	// Create Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JaegerEndpoint)))
	if err != nil {
		return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	// Create resource with service name
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp, sdktrace.WithBatchTimeout(5*time.Second)),
		sdktrace.WithResource(res),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global text map propagator for context propagation
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Telemetry{
		TracerProvider: tp,
		Tracer:         tp.Tracer(cfg.ServiceName),
	}, nil
}

// Shutdown gracefully shuts down the telemetry
func (t *Telemetry) Shutdown(ctx context.Context) error {
	if t.TracerProvider == nil {
		return nil
	}
	return t.TracerProvider.Shutdown(ctx)
}
