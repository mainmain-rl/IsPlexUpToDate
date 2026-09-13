# Building stage
FROM golang:1.25.6-alpine AS builder

WORKDIR /build

COPY go.mod ./
COPY internal/. internal/.
COPY cmd/. cmd/.

RUN go mod download 2>/dev/null || true

# TARGETOS/TARGETARCH sont fournis automatiquement par buildx pour chaque
# plateforme cible listée dans `platforms:` du workflow (amd64, arm64, ...).
ARG TARGETOS
ARG TARGETARCH

# Build the application
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /out/isplexuptodate ./cmd/isplexuptodate


# Final Stage
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/isplexuptodate /isplexuptodate
USER nonroot:nonroot
ENTRYPOINT ["/isplexuptodate"]
