.PHONY: help

.DEFAULT_GOAL := help

build: ## build for elk audio os
	mkdir -p bin
	GOOS=linux go build -o bin/sushi-midi-osc-bridge cmd/main.go

install-go:
	rm -rf /usr/local/go
	wget -O /tmp/go.tar.gz https://go.dev/dl/go1.20.6.linux-amd64.tar.gz
	sudo tar -C /usr/local -xzf /tmp/go.tar.gz
	rm /tmp/go.tar.gz

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'



