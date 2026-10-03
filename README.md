# OFM API Gateway

## Purpose

The API Gateway is the public HTTP entry point for OFM. It validates public requests and forwards them to owning services; it does not own business data, password hashing, or saga state. Status: active.

## Boundaries and interfaces

- Exposes the versioned public HTTP API.
- Calls registration, order, user, gig, file, chat, payment, and review services through internal contracts.
- Owns transport validation, authentication middleware, routing, and configured fallback policy.
- Does not persist business entities or publish service-owned domain events.

Canonical HTTP and gRPC contracts live in the gateway handlers and ofm-common/proto. WebSocket connections are handled by realtime-service.

## Local development

Run from this repository:

    cp .env.example .env
    just run
    go test ./...

Configuration is read from .env: APP_ENV and LOG_LEVEL select runtime behavior, HTTP_* controls the public listener, service *_ADDRESS values select internal gRPC targets, and JWT settings control token validation. Do not commit .env or secrets. Run the complete stack from ofm-infra.

## Build and operations

Dockerfile builds ofm/api-gateway:<tag> and ofm-infra deploys the api-gateway Helm workload. Use structured logs, OpenTelemetry traces, and gateway metrics to diagnose validation and downstream failures. Local ports are defined in PORTS.md.

