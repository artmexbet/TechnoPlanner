FROM golang:1.25.4-alpine3.22 AS builder
ARG service_name
ARG env=stage
ARG cmd_path=cmd/main.go
WORKDIR /app
COPY . .
RUN go build -o app ./services/${service_name}/${cmd_path}

FROM alpine:3.22
ARG service_name
ARG env=stage
WORKDIR /app
COPY --from=builder /app/app ./
CMD ["./app"]
