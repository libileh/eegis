#!/bin/bash

# Define the order of microservices to build
MICROSERVICES=("notifications" "users" "posts")

# Function to build a Docker image for a microservice
build_microservice() {
    local service=$1
    echo "Building Docker image for $service..."
    docker build -t $service:latest -f $service/Dockerfile .
    if [ $? -eq 0 ]; then
        echo "Successfully built Docker image for $service."
    else
        echo "Failed to build Docker image for $service."
        exit 1
    fi
}

# Build Docker images in the specified order
for service in "${MICROSERVICES[@]}"; do
    build_microservice $service
done

echo "All Docker images built successfully."