FROM golang:1.23

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Download dependencies
RUN go mod tidy

# Copy the rest of the application files
COPY . .

# Generate templates
RUN templ generate -path .

# Build the Go application
RUN go build -v -o bin ./cmd/api

EXPOSE 8080

ENTRYPOINT ["/app/bin"]
