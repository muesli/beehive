BEEHIVE_VERSION=0.3.99
COMMIT_SHA=$(shell git rev-parse --short HEAD)

all: submodule embed

submodule:
	git submodule update --init

noembed: submodule build

generate:
	$(shell go env GOPATH)/bin/go-bindata --tags embed --pkg api -o api/bindata.go --ignore config/.git assets/... config/...

go-bindata:
	[ -f $(shell go env GOPATH)/bin/go-bindata ] || go get -u github.com/kevinburke/go-bindata/go-bindata

embed: go-bindata generate build

build:
	go build -tags 'embed' -ldflags '-s -w -X main.Version=$(BEEHIVE_VERSION) -X main.CommitSHA=$(COMMIT_SHA)'

debug: submodule go-bindata generate
	go build -tags 'embed' -ldflags '-X main.Version=$(BEEHIVE_VERSION) -X main.CommitSHA=$(COMMIT_SHA)'

test:
	go test -v $(shell go list ./... | grep -v vendor/)

release:
	@./tools/release.sh

clean:
	rm -f beehive

.PHONY: clean embed go-bindata noembed generate submodule build release all
