.PHONY: help run build docker-build docker-run test test-valid test-invalid health clean

help:
	@echo "Integration Service Demo - Available Commands:"
	@echo ""
	@echo "  make run           - Run the service locally"
	@echo "  make build         - Build the Go binary"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run service in Docker"
	@echo "  make test          - Run all test requests"
	@echo "  make test-valid    - Test with valid token"
	@echo "  make test-invalid  - Test with invalid token"
	@echo "  make health        - Check service health"
	@echo "  make clean         - Clean build artifacts"

run:
	@echo "Starting Integration Service..."
	go run main.go

build:
	@echo "Building binary..."
	go build -o bin/integration-service main.go
	@echo "Binary created at bin/integration-service"

docker-build:
	@echo "Building Docker image..."
	docker build -t integration-demo:latest .

docker-run:
	@echo "Running Docker container..."
	docker run -p 9090:9090 integration-demo:latest

test:
	@echo "Running test suite..."
	@./test-requests.sh

test-valid:
	@echo "Testing with valid token..."
	@curl -s -X POST http://localhost:9090/api/process \
		-H "Content-Type: application/json" \
		-d @examples/valid-request.json | jq '.'

test-invalid:
	@echo "Testing with invalid token..."
	@curl -s -X POST http://localhost:9090/api/process \
		-H "Content-Type: application/json" \
		-d @examples/invalid-token-request.json | jq '.'

health:
	@echo "Checking service health..."
	@curl -s http://localhost:9090/health | jq '.'

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "Clean complete"

