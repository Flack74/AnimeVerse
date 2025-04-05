# debian for easier build utilities
FROM golang:1.24-bullseye AS build-base

WORKDIR /app

# Copy only go mod files first for caching
COPY go.mod go.sum ./

# Cache modules
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# ========================== Development Stage ==========================
FROM build-base AS dev

# Install air & delve
RUN go install github.com/air-verse/air@v1.61.7 && \
    go install github.com/go-delve/delve/cmd/dlv@v1.22.1

WORKDIR /app

# Copy entire app
COPY . .

# Create build output directory
RUN mkdir -p tmp

# Copy Air config if it exists
COPY .air.toml .air.toml

# Use Air to hot-reload
CMD ["air", "-c", ".air.toml"]

# ========================= Production Build Stage =========================
FROM build-base AS build-production

# Add non-root user
RUN useradd -u 1001 nonroot

WORKDIR /app

COPY . .

# Build healthcheck binary
RUN go build \
    -buildvcs=false \
    -ldflags="-linkmode external -extldflags -static" \
    -tags netgo \
    -o /app/healthcheck \
    ./healthcheck/healthcheck.go

# Build main app
RUN go build \
    -buildvcs=false \
    -ldflags="-linkmode external -extldflags -static" \
    -tags netgo \
    -o /app/animeverse-api \
    ./cmd/server

# ========================= Final Minimal Image =========================
FROM scratch

# Set gin to release mode
ENV GIN_MODE=release

WORKDIR /

# Non-root setup
COPY --from=build-production /etc/passwd /etc/passwd
COPY --from=build-production /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy built binaries
COPY --from=build-production /app/healthcheck healthcheck
COPY --from=build-production /app/animeverse-api animeverse-api

USER nonroot

EXPOSE 8000

CMD ["/animeverse-api"]