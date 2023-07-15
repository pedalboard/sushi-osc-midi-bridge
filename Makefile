.PHONY: help

.DEFAULT_GOAL := help

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

