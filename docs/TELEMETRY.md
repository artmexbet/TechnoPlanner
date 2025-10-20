# OpenTelemetry Integration with Jaeger

This document describes the OpenTelemetry integration with Jaeger for collecting metrics and traces.

## Overview

The application now includes OpenTelemetry instrumentation for both the auth and gateway services, with traces being exported to Jaeger.

## Components

### Telemetry Library (`libs/telemetry`)

A shared library that provides OpenTelemetry configuration and initialization with Jaeger exporter.

**Configuration:**
- `service_name`: Name of the service for identification in Jaeger
- `jaeger_endpoint`: URL of the Jaeger collector (default: `http://localhost:14268/api/traces`)
- `enabled`: Boolean to enable/disable telemetry

### Auth Service

The auth service now includes:
- OpenTelemetry tracer initialization
- gRPC server instrumentation with `otelgrpc.NewServerHandler()`
- Automatic trace propagation for all gRPC calls

### Gateway Service

The gateway service now includes:
- OpenTelemetry tracer initialization
- HTTP middleware for tracing Fiber requests
- gRPC client instrumentation for calls to the auth service
- Automatic trace propagation across HTTP and gRPC boundaries

## Running with Jaeger

### 1. Start Jaeger

Use docker-compose to start all services including Jaeger:

```bash
cd deploy
docker-compose up -d jaeger
```

### 2. Access Jaeger UI

Once Jaeger is running, you can access the UI at:
- **Jaeger UI**: http://localhost:16686

### 3. Start the Services

Start the auth and gateway services with telemetry enabled (already enabled in default configs).

The services will automatically send traces to Jaeger.

## Configuration

### Auth Service (`services/auth/cmd/config/cfg.yaml`)

```yaml
telemetry:
  service_name: "auth-service"
  jaeger_endpoint: "http://localhost:14268/api/traces"
  enabled: true
```

### Gateway Service (`services/gateway/cmd/server/config/cfg.yaml`)

```yaml
telemetry:
  service_name: "gateway-service"
  jaeger_endpoint: "http://localhost:14268/api/traces"
  enabled: true
```

## Viewing Traces

1. Open Jaeger UI at http://localhost:16686
2. Select the service name from the dropdown (e.g., `auth-service` or `gateway-service`)
3. Click "Find Traces" to view all traces
4. Click on a specific trace to see detailed span information

## Trace Propagation

Traces are automatically propagated:
- From gateway HTTP requests to gRPC calls to auth service
- Across all internal service calls
- Redis operations in the auth service are also instrumented

## Disabling Telemetry

To disable telemetry, set `enabled: false` in the service configuration:

```yaml
telemetry:
  enabled: false
```

## Available Ports

Jaeger exposes multiple ports:
- `5775/udp`: Zipkin compact thrift protocol (deprecated)
- `6831/udp`: Jaeger compact thrift protocol
- `6832/udp`: Jaeger binary thrift protocol
- `5778`: Jaeger HTTP configurations
- `16686`: Jaeger UI
- `14268`: Jaeger HTTP collector
- `14250`: Jaeger gRPC collector
- `9411`: Zipkin compatible endpoint

## Dependencies

The following OpenTelemetry dependencies have been added:
- `go.opentelemetry.io/otel@v1.38.0`
- `go.opentelemetry.io/otel/exporters/jaeger@v1.17.0`
- `go.opentelemetry.io/otel/sdk@v1.38.0`
- `go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc@v0.58.0`

These dependencies have been checked for vulnerabilities and are safe to use.
