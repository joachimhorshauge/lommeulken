FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -v -o bin ./cmd/api

EXPOSE 8080

ENTRYPOINT ["/app/bin"]
