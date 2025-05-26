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
    libasound2 \
    openssh-client

WORKDIR /root/

# Create SSH directory and proper config
RUN mkdir -p /root/.ssh && \
    echo "Host *\n\tStrictHostKeyChecking no\n\tUserKnownHostsFile /dev/null" > /root/.ssh/config && \
    chmod 600 /root/.ssh/config && \
    chmod 700 /root/.ssh

# Also set the environment variable explicitly
ENV SSH_KNOWN_HOSTS=/root/.ssh/known_hosts

# Copy the binary from builder
COPY --from=builder /app/main .

# Command to run the application
ENTRYPOINT ["bash", "-c", "eval $(ssh-agent -s) && ./main"]
