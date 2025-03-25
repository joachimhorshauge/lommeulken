# Stage 1: Build
FROM golang:1.24-alpine AS build

# Install curl and required libraries for tailwindcss
RUN apk add --no-cache curl libstdc++ libgcc

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install github.com/a-h/templ/cmd/templ@latest && \
    templ generate

# Determine the architecture and download the correct tailwindcss binary
ARG TARGETARCH
RUN if [ "$TARGETARCH" = "arm64" ]; then \
        curl -sL https://github.com/tailwindlabs/tailwindcss/releases/download/v4.0.6/tailwindcss-linux-arm64-musl -o tailwindcss; \
    elif [ "$TARGETARCH" = "amd64" ]; then \
        curl -sL https://github.com/tailwindlabs/tailwindcss/releases/download/v4.0.6/tailwindcss-linux-x64-musl -o tailwindcss; \
    else \
        echo "Unsupported architecture: $TARGETARCH"; exit 1; \
    fi && \
    chmod +x tailwindcss && \
    ls -lah

# Run tailwindcss to generate CSS
RUN ./tailwindcss -i cmd/web/styles/input.css -o cmd/web/assets/css/output.css

# Build the Go application
RUN go build -o main cmd/api/main.go

# Stage 2: Production
FROM alpine:3.20.1 AS prod

# Install required libraries for tailwindcss in the production image (if needed)
RUN apk add --no-cache libstdc++ libgcc

WORKDIR /app
COPY --from=build /app/main /app/main
COPY --from=build /app/cmd/web/assets/css/output.css /app/cmd/web/assets/css/output.css
EXPOSE ${PORT}
CMD ["./main"]
