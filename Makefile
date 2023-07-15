.PHONY: help

.DEFAULT_GOAL := help

generate:
	protoc \
		--proto_path=sushi-grpc-api \
		--go-grpc_out=internal/sushi_rpc \
		--go-grpc_opt=paths=source_relative \
		--go-grpc_opt=Msushi_rpc.proto=github.com/pedalboard/somb/internal/sushi_rpc \
		sushi_rpc.proto

build: ## build
	mkdir -p bin
	go build -o bin/sushi-midi-osc-bridge cmd/sushi_midi_osc_bridge.go

run:
	bin/sushi-midi-osc-bridge

install-go: ## install go
	wget -O /tmp/go.tar.gz https://go.dev/dl/go1.20.6.linux-arm64.tar.gz
	sudo rm -rf /usr/local/go
	sudo tar -C /usr/local -xzf /tmp/go.tar.gz
	rm /tmp/go.tar.gz

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

