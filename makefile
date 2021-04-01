# build file
GOCMD=go
# Use -a flag to prevent code cache problems.
GOBUILD=$(GOCMD) build -ldflags -s -v -a

rpc-mac: BIN_BINARY_NAME=rpc_server
rpc-mac: GOOS=darwin
rpc-mac: GOARCH=amd64
rpc-mac:
	$(GOBUILD) -o $(BIN_BINARY_NAME) cmd/main.go
	mv $(BIN_BINARY_NAME) bin/mac/
	@echo "Build $(BIN_BINARY_NAME) successfully. You can run bin/$(BIN_BINARY_NAME) now.If you can't see it soon,wait some seconds"

rpc-win: BIN_BINARY_NAME=rpc_server.exe
rpc-win: GOOS=windows
rpc-win: GOARCH=amd64
rpc-win:
	$(GOBUILD) -o $(BIN_BINARY_NAME) cmd/main.go
	mv $(BIN_BINARY_NAME) bin/win/
	@echo "Build $(BIN_BINARY_NAME) successfully. You can run bin/$(BIN_BINARY_NAME) now.If you can't see it soon,wait some seconds"

rpc-linux: BIN_BINARY_NAME=rpc_server
rpc-linux: GOOS=linux
rpc-linux: GOARCH=amd64
rpc-linux:
	$(GOBUILD) -o $(BIN_BINARY_NAME) cmd/main.go
	mv $(BIN_BINARY_NAME) bin/linux/
	@echo "Build $(BIN_BINARY_NAME) successfully. You can run bin/$(BIN_BINARY_NAME) now.If you can't see it soon,wait some seconds"
