# Build stage
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/server ./unheicd/main.go

# Final stage
FROM debian:12

# Create a non-root user and set up their home directory
RUN useradd -m -d /home/appuser -s /bin/bash appuser && \
    chown -R appuser:appuser /home/appuser

# Create app directory
RUN mkdir -p /app

# Copy your application binary
COPY --from=builder /app/server /app/server

# Set ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Set the HOME environment variable explicitly
ENV HOME=/home/appuser

# Set the working directory
WORKDIR /app

# Run the application
CMD ["./server"]

# Expose the application port
EXPOSE 8080