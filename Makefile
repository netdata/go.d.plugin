GO  := go

all: vet test build

build:
	mkdir -p dist
	$(GO) build -o dist/godplugin github.com/netdata/go.d.plugin/cmd/godplugin

clean:
	rm -rf dist

test:
	$(GO) test ./... -race -cover -covermode=atomic

vet:
	$(GO) vet ./...