BINARY_NAME=my-website

all: clean format build run

build:
	@echo "Building..."
	@go build -tags netgo -a -v -o bin/$(BINARY_NAME) cmd/main.go

clean:
	@echo "Cleaning..."
	@go clean
	@rm bin/$(BINARY_NAME)

format:
	@echo "Formatting..."
	@go fmt ./...

run:
	@echo "Running..."
	@./bin/$(BINARY_NAME)

test:
	@echo "Testing..."
	@go test ./...
