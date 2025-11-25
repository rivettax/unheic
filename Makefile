.PHONY: dev test help

# Default target
help:
	@echo "Available targets:"
	@echo "  dev     - Start the development server"
	@echo "  test    - Run tests"

# Start the development server
dev:
	@echo "Starting development server..."
	cd unheicd && go run main.go

# Run tests
test:
	@echo "Running tests..."
	cd unheicd && go test ./... 
