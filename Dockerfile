FROM golang:1.23

WORKDIR /app

COPY go.mod .
COPY main.go .

RUN go build -o bin .

EXPOSE 8080

ENTRYPOINT ["/app/bin"]

