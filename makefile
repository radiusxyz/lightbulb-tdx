PROTOC_GEN_GO := $(shell go env GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(shell go env GOPATH)/bin/protoc-gen-go-grpc

PROTO_DIR := proto
OUT_DIR := .
BIN_DIR := bin
SERVER_MAIN := server/main.go
SERVER_OUTPUT := $(BIN_DIR)/server
SERVER_ADDRESS := localhost:50051
CLIENT_MAIN := client/main.go
CLIENT_OUTPUT := $(BIN_DIR)/client

build: build-server build-client

build-server:
	@mkdir -p $(BIN_DIR)
	go build -o $(SERVER_OUTPUT) $(SERVER_MAIN)

build-client:
	@mkdir -p $(BIN_DIR)
	go build -o $(CLIENT_OUTPUT) $(CLIENT_MAIN)

serve: build-server
	$(SERVER_OUTPUT)

run-client: build-client
	$(CLIENT_OUTPUT)

clean:
	rm -rf $(BIN_DIR)

protogen: protoc-auction protoc-attest protoc-benchmark

protoc-auction:
	protoc --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	       $(PROTO_DIR)/auction/auction.proto

protoc-attest:
	protoc --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	       $(PROTO_DIR)/attest/attest.proto

protoc-benchmark:
	protoc --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	       $(PROTO_DIR)/benchmark/benchmark.proto

reflect:
	grpcurl -plaintext $(SERVER_ADDRESS) list

test-rpc:
	./scripts/test_rpc.sh $(rpc)

copy-env:
	@if [ ! -f $(ENV_FILE) ]; then \
		echo "Copying $(ENV_EXAMPLE) to $(ENV_FILE)"; \
		cp $(ENV_EXAMPLE) $(ENV_FILE); \
	else \
		echo "$(ENV_FILE) already exists."; \
	fi

.PHONY: build serve run-client clean protogen reflect test-rpc copy-env