# syntax=docker/dockerfile:1.5

FROM golang:1.24 AS generate-protobuf

# Install system deps with caching
RUN --mount=type=cache,target=/var/cache/apt \
    apt-get update && \
    apt-get install -y --no-install-recommends protobuf-compiler patch ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Install Go protoc plugin
RUN --mount=type=cache,target=/go/pkg/mod \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Copy source files
COPY . /app
WORKDIR /app

# Generate protos
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    cd /tmp && \
    git clone --depth 1 https://github.com/meshtastic/protobufs.git && \
    patch -p1 < /app/patch/remove-nanopb.patch && \
    rm -rf /app/internal/meshtastic/generated && \
    protoc -I=protobufs \
        --go_out=/app/internal/meshtastic \
        --go_opt=module=github.com/meshtastic/go \
        protobufs/meshtastic/*.proto

# Build app
FROM golang:1.24 AS build

WORKDIR /build
COPY --from=generate-protobuf /app /build

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go mod tidy && \
    go build -v -o /usr/bin/meshobserv ./cmd/meshobserv

# Final runtime
FROM debian:bookworm-slim AS runtime
COPY --from=build /usr/bin/meshobserv /usr/bin/meshobserv
RUN \
    apt update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    mkdir -p /certs && \
    chmod 755 /certs
ENTRYPOINT ["/usr/bin/meshobserv"]
