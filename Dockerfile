# --- Build stage ---
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the application
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# --- Runtime stage ---
FROM alpine:3.18

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port your app listens on
EXPOSE 5000

# Run the binary
CMD ["./main"]
