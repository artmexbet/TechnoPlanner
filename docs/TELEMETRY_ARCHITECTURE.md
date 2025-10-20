# OpenTelemetry Integration Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Client Requests                              │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
                    ┌────────────────┐
                    │   Gateway      │
                    │   Service      │
                    │                │
                    │  - HTTP Tracing│◄──── Telemetry Middleware
                    │  - gRPC Client │      (telemetry_middleware.go)
                    └────────┬───────┘
                             │
                             │ Trace Context
                             │ Propagated
                             ▼
                    ┌────────────────┐
                    │   Auth         │
                    │   Service      │
                    │                │
                    │  - gRPC Server │◄──── otelgrpc.NewServerHandler()
                    │  - Redis Ops   │      (already instrumented)
                    └────────┬───────┘
                             │
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
        ▼                    ▼                    ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│   Jaeger     │    │  PostgreSQL  │    │    Redis     │
│  Collector   │    │              │    │              │
│              │    │              │    │  (already    │
│  Port: 14268 │    │              │    │  instrumented)│
└──────┬───────┘    └──────────────┘    └──────────────┘
       │
       │
       ▼
┌──────────────┐
│  Jaeger UI   │
│              │
│  Port: 16686 │
└──────────────┘
```

## Trace Flow

1. **Gateway Service**:
   - Receives HTTP request
   - `TelemetryMiddleware` starts a span
   - Context propagated to gRPC client
   - `otelgrpc.NewClientHandler()` continues the trace

2. **Auth Service**:
   - Receives gRPC request with trace context
   - `otelgrpc.NewServerHandler()` extracts and continues the trace
   - Redis operations are automatically traced (redisotel)
   - Span completed when response sent

3. **Jaeger**:
   - Receives spans from both services
   - Correlates spans into complete traces
   - Displays in UI for analysis

## Configuration Files

- `services/auth/cmd/config/cfg.yaml` - Auth service telemetry config
- `services/gateway/cmd/server/config/cfg.yaml` - Gateway service telemetry config
- `deploy/docker-compose.yml` - Jaeger container configuration

## Key Files

- `libs/telemetry/telemetry.go` - Shared telemetry initialization
- `services/auth/cmd/main.go` - Auth service with gRPC tracing
- `services/gateway/cmd/server/main.go` - Gateway service initialization
- `services/gateway/internal/app/telemetry_middleware.go` - HTTP tracing
- `services/gateway/internal/service/grpcwrapper.go` - gRPC client tracing
