FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main ./main.go

FROM alpine:latest

# Install CA certificates
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Command to run the executable
CMD ["./main"]
