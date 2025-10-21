package server

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor логирует входящие gRPC запросы
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Извлекаем trace ID из контекста
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	slog.InfoContext(ctx, "gRPC request started",
		"method", info.FullMethod,
		"trace_id", traceID,
		"span_id", spanID,
	)

	// Выполняем запрос
	resp, err := handler(ctx, req)

	// Логируем результат
	duration := time.Since(start)
	code := codes.OK
	if err != nil {
		code = status.Code(err)
	}

	logLevel := slog.LevelInfo
	if err != nil {
		logLevel = slog.LevelError
	}

	slog.Log(ctx, logLevel, "gRPC request completed",
		"method", info.FullMethod,
		"duration_ms", duration.Milliseconds(),
		"status", code.String(),
		"trace_id", traceID,
		"span_id", spanID,
		"error", err,
	)

	return resp, err
}
