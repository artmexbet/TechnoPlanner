package opentelemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
)

func newTracerProvider(exp sdktrace.SpanExporter, serviceName string) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

// NewTracerProvider создает TracerProvider с указанным экспортёром и именем сервиса.
// Возвращает провайдер и функцию shutdown, которую необходимо вызвать при завершении работы
// чтобы корректно выгрузить и отправить оставшиеся спаны.
func NewTracerProvider(exp sdktrace.SpanExporter, serviceName string) (*sdktrace.TracerProvider, func(ctx context.Context) error) {
	tp := newTracerProvider(exp, serviceName)

	shutdown := func(ctx context.Context) error {
		// Дадим SDK немного времени на выгрузку батчей перед завершением
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return tp.Shutdown(ctx)
	}

	return tp, shutdown
}

// NewPropagator создает TextMapPropagator для передачи trace context между сервисами.
// Использует W3C Trace Context и W3C Baggage стандарты.
func NewPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newMetricProvider(reader metric.Reader, serviceName string) *metric.MeterProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		),
	)

	if err != nil {
		panic(err)
	}

	return metric.NewMeterProvider(
		metric.WithReader(reader),
		metric.WithResource(r),
	)
}

// NewMetricProvider создает MeterProvider с указанным ридером и именем сервиса.
// Возвращает провайдер и функцию shutdown, которую необходимо вызвать при завершении работы
// чтобы корректно выгрузить и отправить оставшиеся метрики.
func NewMetricProvider(reader metric.Reader, serviceName string) (*metric.MeterProvider, func(ctx context.Context) error) {
	mp := newMetricProvider(reader, serviceName)

	shutdown := func(ctx context.Context) error {
		// Дадим SDK немного времени на выгрузку метрик перед завершением
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return mp.Shutdown(ctx)
	}

	return mp, shutdown
}
