FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY shared ./shared
COPY http-scan ./http-scan
WORKDIR /app/http-scan
RUN CGO_ENABLED=0 GOOS=linux go build -o service .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/http-scan/service .
EXPOSE 9002
CMD ["./service"]
