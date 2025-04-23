FROM golang:1.23 AS builder

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY ./internal ./internal
COPY ./cmd ./cmd

WORKDIR /app/cmd

RUN go build -o /app/main .

CMD ["/app/main"]

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]