# Binary output directory
BIN_DIR := bin

# Binaries
API_BIN := $(BIN_DIR)/jobdoe-api

# Go files
API_MAIN := ./cmd/jobdoe/main.go

# Default target
.PHONY: all
all: build

# Build all binaries
.PHONY: build
build: $(API_BIN)

$(API_BIN):
	@mkdir -p $(BIN_DIR)
	go build -o $(API_BIN) $(API_MAIN)

# Run tests
.PHONY: test
test:
	go test ./...

# Run the API server
.PHONY: run
run:
	go run $(API_MAIN)

# Generate Swagger docs
.PHONY: swagger
swagger:
	./scripts/gen-swagger.sh

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
