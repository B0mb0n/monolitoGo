FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY shared ./shared
COPY tls-scan ./tls-scan
WORKDIR /app/tls-scan
RUN CGO_ENABLED=0 GOOS=linux go build -o service .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/tls-scan/service .
EXPOSE 9003
CMD ["./service"]
