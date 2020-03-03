DEST_DIR=../src/grpc/internal

generate:
	mkdir -p $(DEST_DIR)
	protoc --go_out=plugins=grpc:$(DEST_DIR) *.proto