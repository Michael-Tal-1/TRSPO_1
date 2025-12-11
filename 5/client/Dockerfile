FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY client.go .
RUN go build -o client client.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/client .

CMD ["./client"]
