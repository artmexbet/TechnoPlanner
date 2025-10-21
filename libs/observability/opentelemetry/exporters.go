package opentelemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewOTLPGRPCExporter создает OTLP gRPC exporter.
// endpoint — хост:порт collector'а (по умолчанию "localhost:4317").
// insecureConn — если true, будет использовано небезопасное соединение (без TLS).
func NewOTLPGRPCExporter(ctx context.Context, endpoint string, insecureConn bool) (*otlptrace.Exporter, error) {
	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())), otlptracegrpc.WithInsecure())
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
