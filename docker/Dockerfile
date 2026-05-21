# =============================================================
# devsecops-pipeline-demo — Multi-stage Dockerfile
# =============================================================
# Stage 1: Build  → compila el binario con todas las dependencias
# Stage 2: Runtime → imagen mínima con solo el binario
# =============================================================

# ── Stage 1: Build ───────────────────────────────────────────
FROM golang:1.23-alpine AS builder

# Metadata (OCI standard labels)
LABEL maintainer="Edward Suarez <contact@edwardsuarez.com>"
LABEL org.opencontainers.image.source="https://github.com/edsuarezs/devsecops-pipeline-demo"

# Install CA certificates (needed for HTTPS calls in the binary)
RUN apk add --no-cache ca-certificates

WORKDIR /build

# Copy dependency files first → Docker layer caching
# Only re-downloads modules if go.mod/go.sum change
COPY app/go.mod app/go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY app/ .

# Run tests inside the build (fail fast — don't build broken code)
RUN go test ./... -v

# Compile static binary
# CGO_ENABLED=0 → no C dependencies → runs on scratch/distroless
# -ldflags="-s -w" → strip debug info → smaller binary
# -trimpath → remove local file paths from binary → security
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -trimpath -o /app .

# ── Stage 2: Runtime ─────────────────────────────────────────
FROM gcr.io/distroless/static:nonroot

# Copy CA certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binary
COPY --from=builder /app /app

# Run as non-root user (UID 65532 = nonroot in distroless)
USER nonroot:nonroot

# Expose the application port
EXPOSE 8080

# Use exec form — proper signal handling for graceful shutdown
ENTRYPOINT ["/app"]
