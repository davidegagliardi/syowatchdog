# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o syowatchdog .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Put everything in /app
WORKDIR /app

# Copy binary
COPY --from=builder /app/syowatchdog .

# Create non-root user
RUN adduser -D -s /bin/sh watchdog \
  && chown -R watchdog:watchdog /app

USER watchdog

ENTRYPOINT ["./syowatchdog"]
