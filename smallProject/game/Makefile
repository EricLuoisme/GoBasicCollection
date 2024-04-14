

.PHONY: client server

# 一般来说Makefile都是与go.mod在同一级, 如果不对的情况下, 则指定执行路径
# Variable to store the project root relative to the Makefile
ROOT_DIR := ../

client:
	# Build the client binary
	cd $(ROOT_DIR) && go build -o ./game/bin/client ./game/client/main.go
	# Execute the binary
	./bin/client