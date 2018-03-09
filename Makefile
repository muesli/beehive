BEEHIVE_VERSION=0.2.99
COMMIT_SHA=$(shell git rev-parse --short HEAD)

all: embed

submodule:
	git submodule update --init

noembed: submodule build

generate:
	$(shell go env GOPATH)/bin/go-bindata --tags embed --pkg api -o api/bindata.go --ignore \\.git assets/... config/...

go-bindata:
	[ -f $(shell go env GOPATH)/bin/go-bindata ] || go get -u github.com/jteeuwen/go-bindata/go-bindata

embed: submodule go-bindata generate build

build:
	go build -tags 'embed' -ldflags '-s -w -X main.Version=$(BEEHIVE_VERSION) -X main.CommitSHA=$(COMMIT_SHA)'

debug: submodule go-bindata generate
	go build -tags 'embed' -ldflags '-X main.Version=$(BEEHIVE_VERSION) -X main.CommitSHA=$(COMMIT_SHA)'

test:
	go test -v $(shell go list ./... | grep -v vendor/)

get-deps:
	dep ensure

clean:
	rm -f beehive
.PHONY: clean embed go-bindata get-deps noembed generate submodule build all
