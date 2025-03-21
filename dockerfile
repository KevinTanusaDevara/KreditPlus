FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o kreditplus-app

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/kreditplus-app .

COPY .env .env

EXPOSE 8080

CMD ["./kreditplus-app"]
