# build file
GOCMD=go
# Use -a flag to prevent code cache problems.
GOBUILD=$(GOCMD) build -ldflags -s -v -a

rpc-mac: SERVER_NAME=rpc_server
rpc-mac: CLI_NAME=data_cli
rpc-mac: export GOOS=darwin
rpc-mac: export GOARCH=amd64
rpc-mac:
	$(GOBUILD) -o $(SERVER_NAME) cmd/main.go
	mv $(SERVER_NAME) bin/mac/
	$(GOBUILD) -o $(CLI_NAME) cmd/cli/main.go
	mv $(CLI_NAME) bin/mac/
	@echo "Build $(SERVER_NAME) successfully. You can run bin/$(SERVER_NAME) now.If you can't see it soon,wait some seconds"

rpc-win: SERVER_NAME=rpc_server.exe
rpc-win: CLI_NAME=data_cli.exe
rpc-win: export GOOS=windows
rpc-win: export GOARCH=amd64
rpc-win:
	$(GOBUILD) -o $(SERVER_NAME) cmd/main.go
	mv $(SERVER_NAME) bin/win/
	$(GOBUILD) -o $(CLI_NAME) cmd/cli/main.go
	mv $(CLI_NAME) bin/win/
	@echo "Build $(SERVER_NAME) successfully. You can run bin/$(SERVER_NAME) now.If you can't see it soon,wait some seconds"

rpc-linux: SERVER_NAME=rpc_server
rpc-linux: CLI_NAME=data_cli
rpc-linux: export GOOS=linux
rpc-linux: export GOARCH=amd64
rpc-linux:
	$(GOBUILD) -o $(SERVER_NAME) cmd/main.go
	mv $(SERVER_NAME) bin/linux/
	#$(GOBUILD) -o $(CLI_NAME) cmd/cli/main.go
	#mv $(CLI_NAME) bin/linux/
	@echo "Build $(SERVER_NAME) successfully. You can run bin/$(SERVER_NAME) now.If you can't see it soon,wait some seconds"
