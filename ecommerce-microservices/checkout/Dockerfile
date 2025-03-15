FROM golang:1.23 as builder

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o checkout-service ./checkout/cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/checkout-service .

EXPOSE 8080

CMD ["./checkout-service"]