FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/gateway .
COPY services.json .
EXPOSE 8080
CMD ["./gateway"]
