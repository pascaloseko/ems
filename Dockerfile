# Step 1: Build stage
FROM golang:1.22-alpine as builder

ENV TZ=UTC \
    PATH=/root/go/bin:$PATH

WORKDIR /app

# Install additional dependencies
RUN apk add --no-cache git

# Copy the Go module files and download dependencies
COPY go.mod go.sum /app/
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN go build -o /app/server

### 
## Step 2: Runtime stage
FROM alpine:latest

ENV TZ=UTC \
    PATH=/usr/local/go/bin:$PATH

# Install necessary dependencies
RUN apk add --no-cache ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /app/server /usr/local/bin/server

# Set the entry point for the container
CMD ["server"]
