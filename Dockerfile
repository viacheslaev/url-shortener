# ===========================
# Build stage
# ===========================
# Use a minimal Go image based on Alpine Linux to compile the binary
FROM golang:1.25-alpine AS build

# All source code will live in /src inside the build container
WORKDIR /src

# Install system dependencies required for:
# - fetching Go modules (git)
# - making HTTPS calls (CA certificates)
RUN apk add --no-cache git ca-certificates

# Copy only go.mod and go.sum first.
# This allows Docker to cache dependencies and avoid re-downloading them
# on every build if only source code changes.
COPY go.mod go.sum ./
RUN go mod download

# Now copy the full source code
COPY . .

# Build a static Linux binary for amd64 architecture.
# CGO_ENABLED=0  -> fully static binary (no libc dependency)
# GOOS=linux     -> target OS
# GOARCH=amd64   -> target CPU architecture
# -trimpath      -> removes local file paths from the binary
# -s -w          -> strips debug symbols (smaller binary)
# -o             -> output binary path
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o /out/url-shortener ./cmd/app


# ===========================
# Runtime stage
# ===========================
# Use a tiny Alpine image for the final runtime container
FROM alpine:3.20

# All runtime files will live in /app
WORKDIR /app

# Create a non-root user for security
# Install CA certificates so the app can make HTTPS requests
RUN adduser -D -H -u 10001 appuser && apk add --no-cache ca-certificates

# Copy only the compiled binary from the build stage
COPY --from=build /out/url-shortener /app/url-shortener

# Copy runtime assets:
# - migrations are needed by the migrate container and sometimes by the API
# - openapi is needed by Swagger UI to serve API docs
COPY migrations /app/migrations
COPY openapi /app/openapi

# Drop root privileges and run as an unprivileged user
USER appuser

# Document the HTTP port used by the service
EXPOSE 8080

# Start the API server
ENTRYPOINT ["/app/url-shortener"]
