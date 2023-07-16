.PHONY: help

.DEFAULT_GOAL := help

generate:
	protoc \
		--proto_path=sushi-grpc-api \
		--go_out=internal/sushi_rpc \
		--go_opt=paths=source_relative \
		--go_opt=Msushi_rpc.proto=github.com/pedalboard/somb/internal/sushi_rpc \
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

install: ## install the services into the local system
	$(MAKE) disable-ro
	sudo cp sushi-osc-midi-bridge.service /lib/systemd/system/
	sudo systemctl daemon-reload
	sudo systemctl enable sushi-osc-midi-bridge
	$(MAKE) enable-ro

enable-ro: ## enable overlay fs
	sudo elk_system_utils --remount-as-ro


disable-ro: ## enable overlay fs
	sudo elk_system_utils  --remount-as-rw

status: ## show the service status
	systemctl status sushi-osc-midi-bridge

stop: ## stop the services
	sudo systemctl stop sushi-osc-midi-bridge

start: ## start the services
	sudo systemctl start sushi-osc-midi-bridge

release:
	gh release create --latest --generate-notes $$(git describe --tags --abbrev=0) ./bin/sushi-midi-osc-bridge

install-latest:
	curl -L https://github.com/pedalboard/sushi-osc-midi-bridge/releases/latest/download/sushi-midi-osc-bridge -o bin/sushi-midi-osc-bridge
	chmod +x bin/sushi-midi-osc-bridge

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

