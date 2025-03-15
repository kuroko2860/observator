FROM golang:1.23 as builder

WORKDIR /app

# Install protoc and the Go plugin
RUN apt-get update && apt-get install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

COPY . .

# Generate protobuf code
RUN mkdir -p inventory/proto
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    inventory/proto/inventory.proto

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o inventory-service ./inventory/cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/inventory-service .

EXPOSE 8083

CMD ["./inventory-service"]