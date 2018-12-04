GO  := go

all: download vet test build

download:
	go mod download

build:
	mkdir -p dist
	$(GO) build -o dist/godplugin github.com/netdata/go.d.plugin/cmd/godplugin

clean:
	rm -rf dist

test:
	$(GO) test ./... -race -cover -covermode=atomic

vet:
	$(GO) vet ./...

dev: dev-build dev-up

dev-build:
	docker-compose build

dev-up:
	docker-compose up -d

dev-exec:
	docker-compose exec netdata bash

dev-log:
	docker-compose logs -f netdata