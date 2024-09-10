# Stage 1: Build the Go binary
FROM golang:1.20-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary
RUN go build -o server .

# Stage 2: Create a lightweight image with the binary
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server .

# Copy the static directory
COPY static ./static

# Expose the port that the server will listen on
EXPOSE 8080

# Run the server
CMD ["./server"]

