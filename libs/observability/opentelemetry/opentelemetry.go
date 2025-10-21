package opentelemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
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
