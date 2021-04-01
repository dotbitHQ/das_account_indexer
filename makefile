# build file
GOCMD=go
# Use -a flag to prevent code cache problems.
GOBUILD=$(GOCMD) build -ldflags -s -v -a

rpc: BIN_BINARY_NAME=rpc_server
rpc:
	$(GOBUILD) -o $(BIN_BINARY_NAME) cmd/main.go
	mv $(BIN_BINARY_NAME) bin/
	@echo "Build $(BIN_BINARY_NAME) successfully. You can run bin/$(BIN_BINARY_NAME) now.If you can't see it soon,wait some seconds"

