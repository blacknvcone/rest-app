# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install make and git (if needed for dependencies)
RUN apk add --no-cache make git

# Copy go.mod, go.sum, and Makefile first for better caching
COPY go.mod go.sum Makefile ./

# Install dependencies using Makefile target
RUN make instal-deps

# Copy the rest of the source code
COPY . .

# Build the Go app (replace 'main.go' with your entrypoint if different)
RUN go build -o app

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Expose port (change if your app uses a different port)
EXPOSE 8080

# Run the app
ENTRYPOINT ["./app"]
