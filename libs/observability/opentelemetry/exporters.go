package opentelemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewOTLPGRPCExporter создает OTLP gRPC exporter.
// endpoint — хост:порт collector'а (по умолчанию "localhost:4317").
// insecureConn — если true, будет использовано небезопасное соединение (без TLS).
func NewOTLPGRPCExporter(ctx context.Context, endpoint string, insecureConn bool) (*otlptrace.Exporter, error) {
	if endpoint == "" {
		endpoint = "localhost:4317"
	}

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(endpoint),
	}
	if insecureConn {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exp, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP gRPC exporter: %w", err)
	}
	return exp, nil
}

// NewOTLPHTTPExporter создает OTLP HTTP exporter.
// endpoint — хост:порт collector'а (по умолчанию "localhost:4318").
// insecureConn — если true, будет использовано небезопасное соединение (HTTP вместо HTTPS).
func NewOTLPHTTPExporter(ctx context.Context, endpoint string, insecureConn bool) (sdktrace.SpanExporter, error) {
	if endpoint == "" {
		endpoint = "localhost:4318"
	}

	opts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(endpoint)}
	if insecureConn {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	exp, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP HTTP exporter: %w", err)
	}
	return exp, nil
}

// NewOTLPGRPCMetricReader создает OTLP gRPC metric reader с периодической отправкой метрик.
// endpoint — хост:порт collector'а (по умолчанию "localhost:4317").
// insecureConn — если true, будет использовано небезопасное соединение (без TLS).
// interval — интервал отправки метрик (по умолчанию 10 секунд).
func NewOTLPGRPCMetricReader(ctx context.Context, endpoint string, insecureConn bool, interval time.Duration) (metric.Reader, error) {
	if endpoint == "" {
		endpoint = "localhost:4317"
	}

	if interval == 0 {
		interval = 10 * time.Second
	}

	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(endpoint),
	}
	if insecureConn {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	exporter, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP gRPC metric exporter: %w", err)
	}

	reader := metric.NewPeriodicReader(exporter, metric.WithInterval(interval))
	return reader, nil
}

// NewOTLPHTTPMetricReader создает OTLP HTTP metric reader с периодической отправкой метрик.
// endpoint — хост:порт collector'а (по умолчанию "localhost:4318").
// insecureConn — если true, будет использовано небезопасное соединение (HTTP вместо HTTPS).
// interval — интервал отправки метрик (по умолчанию 10 секунд).
func NewOTLPHTTPMetricReader(ctx context.Context, endpoint string, insecureConn bool, interval time.Duration) (metric.Reader, error) {
	if endpoint == "" {
		endpoint = "localhost:4318"
	}

	if interval == 0 {
		interval = 10 * time.Second
	}

	opts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(endpoint),
	}
	if insecureConn {
		opts = append(opts, otlpmetrichttp.WithInsecure())
	}

	exporter, err := otlpmetrichttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP HTTP metric exporter: %w", err)
	}

	reader := metric.NewPeriodicReader(exporter, metric.WithInterval(interval))
	return reader, nil
}
