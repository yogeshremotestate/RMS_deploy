# Step 1: Build stage
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./


# Download dependencies
RUN go mod tidy && go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Step 2: Final image (using a more recent Debian-based image)
FROM debian:bookworm-slim

# Set the working directory for the final image
WORKDIR /root/

# Install necessary dependencies: bash, curl, libc6 (for GLIBC compatibility)
RUN apt-get update && apt-get install -y ca-certificates bash curl libc6

# Copy the built application and migrations from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations


# Expose the port the app will run on
EXPOSE 8080

# Run the application
CMD ["./main"]





