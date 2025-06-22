# syntax=docker/dockerfile:1.4

# ========================= Base Stage =========================
FROM golang:1.24-bullseye AS build-base
WORKDIR /src

# only copy modules, download deps
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# ========================= Development Stage =========================
FROM build-base AS dev

# install hot‑reload & debugger
RUN go install github.com/air-verse/air@v1.61.7 \
 && go install github.com/go-delve/delve/cmd/dlv@v1.22.1

WORKDIR /src
COPY . .
RUN mkdir -p tmp

CMD ["air", "-c", ".air.toml"]

# ========================= Production Build Stage =========================
FROM build-base AS build-production

LABEL \
  org.opencontainers.image.title="Animeverse API" \
  org.opencontainers.image.description="REST API for Animeverse" \
  org.opencontainers.image.source="https://github.com/Flack74/animeverse" \
  org.opencontainers.image.license="MIT"

# create non‑root user early
RUN useradd -u 1001 nonroot

# install CA certs for TLS
RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /src
ENV CGO_ENABLED=0 \
    GIN_MODE=release

# copy code & build
COPY . .
RUN go build \
    -buildvcs=false \
    -tags netgo \
    -ldflags="-s -w" \
    -o /animeverse-api \
    .

# ========================= Final Minimal Image =========================
FROM scratch

ENV GIN_MODE=release
WORKDIR /

# bring in user & certs
COPY --from=build-production /etc/passwd /etc/passwd
COPY --from=build-production /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# copy the API binary
COPY --from=build-production /animeverse-api /animeverse-api

USER nonroot

EXPOSE 8000
ENTRYPOINT ["/animeverse-api"]