PROJECT=tc-proxy
DEST_DIR=src/grpc/internal
SRC_DIR=proto

generate:
	mkdir -p $(DEST_DIR)
	protoc -I $(SRC_DIR) --go_out=plugins=grpc:$(DEST_DIR) $(SRC_DIR)/proxy.proto 

build:
	docker build --tag $(PROJECT) .
