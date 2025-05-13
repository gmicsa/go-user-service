# ---- Build Stage ----
FROM golang:1.24.1-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download # Download dependencies

COPY . .

ARG APP_VERSION="dev"
RUN GIT_COMMIT=$(git rev-parse --short HEAD || echo "unknown") && \
    BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') && \
    CGO_ENABLED=0 go build \
    -ldflags="-w -s -X 'user-service/health.Version=${APP_VERSION}' -X 'user-service/health.Commit=${GIT_COMMIT}' -X 'user-service/health.BuildDate=${BUILD_DATE}'" \
    -o /app/service .

# ---- Final Stage ----
FROM scratch AS final

# Add a new user and group (hardcoded UID/GID)
# Use an existing system UID (e.g., 1000) to avoid conflicts
USER 1000:1000

WORKDIR /app

# Copy necessary CA certificates (needed for HTTPS requests from the app)
COPY --from=builder --chown=1000:1000 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder --chown=1000:1000 /app/service /app/service

EXPOSE 8080 8081 8082

ENV ENABLE_PPROF=false
ENV GOMEMLIMIT="512MB"

ENTRYPOINT ["/app/service"]