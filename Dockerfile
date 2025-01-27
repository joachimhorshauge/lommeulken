# Stage 1: Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum files and install dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy the rest of the application files and generate templates
COPY . .
RUN templ generate -path .

# Build the application binary
RUN go build -v -o bin ./cmd/api

# Stage 2: Runtime stage
FROM gcr.io/distroless/static-debian11

WORKDIR /app

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/bin /app/bin

# Expose the application port
EXPOSE 8080

# Set the binary as the entrypoint
ENTRYPOINT ["/app/bin"]
