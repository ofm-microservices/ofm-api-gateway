FROM golang:1.25.5 AS builder

WORKDIR /src

COPY ofm-common /src/ofm-common
COPY ofm-api-gateway /src/ofm-api-gateway

WORKDIR /src/ofm-api-gateway

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api-gateway ./cmd/api-gateway

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && groupadd --system app \
    && useradd --system --gid app --home-dir /app --shell /usr/sbin/nologin app \
    && chown app:app /app \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/api-gateway /app/api-gateway

USER app:app

EXPOSE 8080

CMD ["/app/api-gateway"]
