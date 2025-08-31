.PHONY: build run clean test html

build:
	go build -o ./src/... ./test/...

run:
	go run ./src/...

test:
	go test -v ./tests/... -covermode=count -coverprofile=coverage.out -coverpkg=./src/...

html: test
	go tool cover -html=coverage.out -o coverage.html
	