# Stage 1: Build the application
FROM golang:1.24-bookworm AS builder

# Dependencies
RUN apt update && apt-get install -y \
    alsa-utils \
    libasound2-dev \
    pkg-config


# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY ../go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY .. .

# Build the application
RUN GOOS=linux go build -o main .

# Stage 2: Create the final image
FROM debian:bookworm

# Dependencies
RUN apt update && apt-get install -y \
    alsa-utils \
    libasound2

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Command to run the application
CMD ["./main"]
