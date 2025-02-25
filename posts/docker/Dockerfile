# Base image with Go compiler and Alpine OS
FROM golang:1.23.2-alpine3.19 AS builder

# Set environment variables
ENV GO_VERSION 1.23.2
ENV GOOS linux
ENV CGO_ENABLED 0

# Set working directory
WORKDIR /eegis

# Copy Go module files from the posts directory
COPY posts/go.mod posts/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code from the posts directory
COPY posts/ ./

# Build the application
RUN go build -o posts -v ./cmd/server

# Final image with minimal Alpine OS
FROM alpine:3.19.0

# Set environment variables
ENV PORT 9103

# Install required OS packages
RUN apk --no-cache add \
    ca-certificates \
    curl \
    net-tools \
    tzdata

# Create non-root user and group
RUN addgroup -S appuser && adduser -S appuser -G appuser

# Set working directory
WORKDIR /eegis

# Copy binary from builder
COPY --from=builder /eegis/posts .

# Expose port
EXPOSE ${PORT}

# Run as non-root user
USER appuser:appuser

# Start the application
CMD ["./posts"]