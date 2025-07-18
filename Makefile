.PHONY: dev test help

# Default target
help:
	@echo "Available targets:"
	@echo "  dev     - Start the development server"
	@echo "  test    - Run tests"

# Start the development server
dev:
	@echo "Starting development server..."
	go run cmd/unheicd/main.go

# Run tests
test:
	@echo "Running tests..."
	go test ./... 
