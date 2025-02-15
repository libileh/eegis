# Base image for building the Go application
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY ./backend/go.mod ./
COPY ./backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY ./backend .

# Navigate to the directory containing main.go and build the binary
WORKDIR /app/cmd/server
RUN go build -o /eegis

# Use a minimal image for running the application
FROM alpine:latest

# Set working directory in the container
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /eegis .

# Install timezone data
RUN apk add --no-cache tzdata

# Expose the application port
EXPOSE 9100

# Command to run the application
CMD ["./eegis"]