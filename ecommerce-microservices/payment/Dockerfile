FROM golang:1.23 as builder

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o payment-service ./payment/cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/payment-service .

EXPOSE 8082

CMD ["./payment-service"]