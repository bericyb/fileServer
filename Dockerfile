# Start from the official Golang image for building the application
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go app
RUN go build -o /go-server

# Start a new stage from a smaller image for production (multi-stage build)
FROM alpine:3.18

# Set the working directory inside the container
WORKDIR /app

# Copy the Go app binary from the builder stage
COPY --from=builder /go-server .

# Expose the port that the server will run on
EXPOSE 8080

# Create an 'uploads' directory to store uploaded files
RUN mkdir -p uploads

# Command to run the Go app
CMD ["/app/go-server"]
