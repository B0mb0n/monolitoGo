FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY shared ./shared
COPY port-scan ./port-scan
WORKDIR /app/port-scan
RUN CGO_ENABLED=0 GOOS=linux go build -o service .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/port-scan/service .
EXPOSE 9001
CMD ["./service"]
