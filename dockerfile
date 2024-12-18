# Stage 1: Build the Go binaries
FROM golang:1.23-alpine AS builder 

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files FIRST
COPY go.mod go.sum ./

# Download dependencies and tidy the go.mod file
RUN go mod tidy

# Copy the source code AFTER downloading dependencies
COPY . .

# Build the server binary
RUN go build -o server ./cmd/server

# Build the client binary
RUN go build -o client ./cmd/client

# Stage 2: Create the final image
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binaries from the builder stage
COPY --from=builder /app/server /app/server
COPY --from=builder /app/client /app/store

# Expose the server port
EXPOSE 8080

# Set the entrypoint for the server
CMD ["./server"]