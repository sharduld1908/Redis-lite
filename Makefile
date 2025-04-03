# Makefile for the TCP Server/Client App

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_DIR=./bin
SERVER_BINARY_NAME=server
CLIENT_BINARY_NAME=client
SERVER_CMD_PATH=./cmd/server
CLIENT_CMD_PATH=./cmd/client

# Default target executed when you just run 'make'
all: build

# Build targets
build: build-server build-client

build-server:
	@echo "Building server..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_DIR)/$(SERVER_BINARY_NAME) $(SERVER_CMD_PATH)
	@echo "Server built -> $(BINARY_DIR)/$(SERVER_BINARY_NAME)"

build-client:
	@echo "Building client..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_DIR)/$(CLIENT_BINARY_NAME) $(CLIENT_CMD_PATH)
	@echo "Client built -> $(BINARY_DIR)/$(CLIENT_BINARY_NAME)"

# Run targets (builds first if necessary)
run-server: build-server
	@echo "Running server..."
	@$(BINARY_DIR)/$(SERVER_BINARY_NAME)

run-client: build-client
	@echo "Running client..."
	@$(BINARY_DIR)/$(CLIENT_BINARY_NAME)

# Clean target
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)
	@echo "Cleaned."

# Help target
help:
	@echo "Available commands:"
	@echo "  make build          Build both server and client binaries"
	@echo "  make build-server   Build only the server binary"
	@echo "  make build-client   Build only the client binary"
	@echo "  make run-server     Build (if needed) and run the server"
	@echo "  make run-client     Build (if needed) and run the client"
	@echo "  make clean          Remove build artifacts"
	@echo "  make help           Show this help message"

# Phony targets prevent conflicts with files of the same name
.PHONY: all build build-server build-client run-server run-client clean help