GO  := go
DEP := dep

all: test build

build:
	mkdir -p pkg
	$(GO) build -o pkg/godplugin github.com/netdata/go.d.plugin/cmd/godplugin

clean:
	rm -rf pkg

test:
	$(GO) list ./... | xargs -n1 -I% $(GO) test % -race

cover:
	$(GO) list ./... | xargs -n1 -I% $(GO) test % -cover

dep:
	$(DEP) ensure -v
