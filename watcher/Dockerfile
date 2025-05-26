# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY ../go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY .. .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Create the final image
FROM alpine:latest

# Add CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Command to run the application
CMD ["./main"]
