FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY server.go .
RUN go build -o server server.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 9000

CMD ["./server"]
