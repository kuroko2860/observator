FROM golang:1.23 as builder

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o address-service ./address/cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/address-service .

EXPOSE 8084

CMD ["./address-service"]