module github.com/artmexbet/TechnoPlanner/services/gateway

go 1.25.1

require (
	github.com/artmexbet/TechnoPlanner/libs/broker v0.0.0-20260211113508-95b9b554ae2e
	github.com/artmexbet/TechnoPlanner/libs/config v0.0.0-20260211113508-95b9b554ae2e
	github.com/artmexbet/TechnoPlanner/libs/dto v0.0.0-00010101000000-000000000000
	github.com/artmexbet/TechnoPlanner/libs/observability v0.0.0-20260211113508-95b9b554ae2e
	github.com/artmexbet/TechnoPlanner/libs/proto v0.0.0-20260211113508-95b9b554ae2e
	github.com/go-playground/validator/v10 v10.30.1
	github.com/gofiber/contrib/v3/otel v1.0.0
	github.com/gofiber/fiber/v3 v3.0.0
	github.com/google/uuid v1.6.0
	github.com/mailru/easyjson v0.9.1
	github.com/nats-io/nats.go v1.48.0
	github.com/samber/slog-fiber v1.21.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.65.0
	go.opentelemetry.io/otel v1.41.0
	go.opentelemetry.io/otel/trace v1.41.0
	google.golang.org/grpc v1.78.0
)

require (
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.13 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/gofiber/schema v1.7.0 // indirect
	github.com/gofiber/utils/v2 v2.0.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3 // indirect
	github.com/ilyakaznacheev/cleanenv v1.5.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/nats-io/nkeys v0.4.15 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/tinylib/msgp v1.6.3 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.69.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib v1.40.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.38.0 // indirect
	go.opentelemetry.io/otel/metric v1.41.0 // indirect
	go.opentelemetry.io/otel/sdk v1.40.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.40.0 // indirect
	go.opentelemetry.io/proto/otlp v1.7.1 // indirect
	golang.org/x/crypto v0.48.0 // indirect
	golang.org/x/net v0.50.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260209200024-4cfbd4190f57 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace (
	github.com/artmexbet/TechnoPlanner/libs/broker => ../../libs/broker
	github.com/artmexbet/TechnoPlanner/libs/config => ../../libs/config
	github.com/artmexbet/TechnoPlanner/libs/dto => ../../libs/dto
	github.com/artmexbet/TechnoPlanner/libs/observability => ../../libs/observability
	github.com/artmexbet/TechnoPlanner/libs/proto => ../../libs/proto
	github.com/artmexbet/TechnoPlanner/libs/utills => ../../libs/utills
)
