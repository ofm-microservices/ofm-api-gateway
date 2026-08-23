# OFM API Gateway

## Purpose

`ofm-api-gateway` is the public entrypoint into the OFM microservice system.
It owns HTTP request validation and transport concerns, then forwards work to
internal services. It does not own business persistence and it should not
implement downstream domain logic such as password hashing or registration
orchestration.

Current responsibilities:

- expose public HTTP endpoints
- validate signup payloads
- call `registration-saga-service` over gRPC

## Run

Local process:

```bash
cp .env.example .env
just run
```

Direct Go command:

```bash
set -a && source .env && set +a && go run ./cmd/api-gateway
```

Docker stack from the shared infra repo:

```bash
cd ../ofm-infra
just infra-up
```

## Environment

The service expects the following variables.

```env
APP_ENV=local
LOG_LEVEL=info

HTTP_HOST=0.0.0.0
HTTP_PORT=8080

REGISTRATION_SAGA_ADDRESS=127.0.0.1:9500
```

Notes:

- `REGISTRATION_SAGA_ADDRESS` is the internal gRPC target for registration
  startup.
- Registration and order startup use their owning services' gRPC boundaries;
  the gateway does not connect to a message broker.

## Technologies

Core runtime:

- Go
- Fiber v2 for HTTP transport
- gRPC client for internal service calls
- Uber Fx for dependency wiring
- Zap for structured JSON logging
- `caarlos0/env` + `godotenv` for configuration loading

Main libraries from `go.mod`:

- `github.com/gofiber/fiber/v2`
- `google.golang.org/grpc`
- `go.uber.org/fx`
- `go.uber.org/zap`
- `github.com/caarlos0/env/v11`

## Architecture Notes

- `internal/domain` contains gateway request and response contracts
- `internal/application` performs validation and delegates to downstream clients
- `internal/presentation/http` owns HTTP handlers
- `internal/presentation/grpc` owns the registration saga client
- `internal/fx` wires the process explicitly

The gateway should remain thin. If a behavior belongs to auth, user, mail, or
registration orchestration, keep it there.
