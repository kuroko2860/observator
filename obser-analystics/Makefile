hello:
	echo "Hello"

build:
	go build -o bin/main main/main.go

run:
	swag init -g main/main.go
	go run main/main.go

.PHONY: hello build run
