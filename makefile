PROTOC_GEN_GO := $(shell go env GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(shell go env GOPATH)/bin/protoc-gen-go-grpc

PROTO_DIR := proto
OUT_DIR := .
BIN_DIR := bin
MAIN := server/main.go
OUTPUT := $(BIN_DIR)/lightbulb-tdx

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(OUTPUT) $(MAIN)

serve: build
	$(OUTPUT)

clean:
	rm -rf $(BIN_DIR)

protogen: protoc-auction protoc-attest

protoc-auction:
	protoc --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	       $(PROTO_DIR)/auction/auction.proto

protoc-attest:
	protoc --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	       $(PROTO_DIR)/attest/attest.proto

test_rpc:
	./scripts/test_rpc.sh $(rpc)

.PHONY: build serve clean protogen test_rpc