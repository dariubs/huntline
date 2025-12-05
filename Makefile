.PHONY: help build build-receiver build-server build-migrate run-receiver run-server migrate clean install deps test

# Variables
BINARY_DIR := bin
RECEIVER_DIR := app/main/receiver
SERVER_DIR := app/main/huntline
MIGRATE_DIR := app/main/migrate
RECEIVER_BINARY := $(BINARY_DIR)/receiver
SERVER_BINARY := $(BINARY_DIR)/server
MIGRATE_BINARY := $(BINARY_DIR)/migrate

# Default target
help:
	@echo "HuntLine Makefile Commands:"
	@echo ""
	@echo "  make build          - Build all binaries (receiver, server, migrate)"
	@echo "  make build-receiver - Build the receiver binary"
	@echo "  make build-server   - Build the web server binary"
	@echo "  make build-migrate  - Build the migrate binary"
	@echo ""
	@echo "  make run-receiver   - Run receiver to fetch yesterday's ProductHunt data"
	@echo "  make run-server     - Run the web server"
	@echo "  make migrate        - Run database migrations"
	@echo ""
	@echo "  make receiver-date DATE=2025-01-15  - Fetch data for specific date"
	@echo "  make receiver-repeat                 - Run receiver with daily schedule"
	@echo "  make receiver-historical              - Backfill historical data"
	@echo "  make receiver-last-month              - Update all data from last month"
	@echo ""
	@echo "  make install        - Install Go dependencies"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make deps           - Download dependencies"

# Create bin directory if it doesn't exist
$(BINARY_DIR):
	@mkdir -p $(BINARY_DIR)

# Build all binaries
build: build-receiver build-server build-migrate

# Build receiver binary
build-receiver: $(BINARY_DIR)
	@echo "Building receiver..."
	@go build -o $(RECEIVER_BINARY) ./$(RECEIVER_DIR)
	@echo "Receiver built: $(RECEIVER_BINARY)"

# Build server binary
build-server: $(BINARY_DIR)
	@echo "Building server..."
	@go build -o $(SERVER_BINARY) ./$(SERVER_DIR)
	@echo "Server built: $(SERVER_BINARY)"

# Build migrate binary
build-migrate: $(BINARY_DIR)
	@echo "Building migrate..."
	@go build -o $(MIGRATE_BINARY) ./$(MIGRATE_DIR)
	@echo "Migrate built: $(MIGRATE_BINARY)"

# Run receiver (fetches yesterday's data by default)
run-receiver: build-receiver
	@echo "Running receiver..."
	@$(RECEIVER_BINARY)

# Run receiver for specific date
receiver-date: build-receiver
	@if [ -z "$(DATE)" ]; then \
		echo "Error: DATE is required. Usage: make receiver-date DATE=2025-01-15"; \
		exit 1; \
	fi
	@echo "Running receiver for date: $(DATE)"
	@$(RECEIVER_BINARY) -date $(DATE)

# Run receiver with daily schedule
receiver-repeat: build-receiver
	@echo "Running receiver with daily schedule..."
	@$(RECEIVER_BINARY) -repeat=true

# Run receiver with custom schedule
receiver-schedule: build-receiver
	@if [ -z "$(TIME)" ]; then \
		echo "Error: TIME is required. Usage: make receiver-schedule TIME=13:00"; \
		exit 1; \
	fi
	@echo "Running receiver with schedule: $(TIME)"
	@$(RECEIVER_BINARY) -repeat=true -schedule $(TIME)

# Run historical backfill
receiver-historical: build-receiver
	@echo "Running historical backfill (this may take a while)..."
	@$(RECEIVER_BINARY) -historical=true

# Run receiver to update last month's data
receiver-last-month: build-receiver
	@echo "Running receiver to update last month's data..."
	@$(RECEIVER_BINARY) -last-month=true

# Run server
run-server: build-server
	@echo "Running server..."
	@$(SERVER_BINARY)

# Run migrations
migrate: build-migrate
	@echo "Running migrations..."
	@$(MIGRATE_BINARY)

# Install dependencies
install: deps

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BINARY_DIR)
	@echo "Clean complete"

# Run tests (if you add tests later)
test:
	@echo "Running tests..."
	@go test ./...

