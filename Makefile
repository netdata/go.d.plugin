GO  := go
DEP := dep

all: test build

build:
	mkdir -p pkg
	$(GO) build -o pkg/godplugin github.com/l2isbad/go.d.plugin/cmd/godplugin

clean:
	rm -rf pkg

test:
	$(GO) test ./... -race

cover:
	$(GO) test ./... -cover

dep:
	$(DEP) ensure -v
