all: test embed

submodule:
	git submodule update --init

noembed: submodule build

generate:
	go-bindata --pkg api -o api/bindata.go --ignore \\.git assets/... config/...

go-bindata:
	[ -f $(shell go env GOPATH)/bin/go-bindata ] || go get -u github.com/jteeuwen/go-bindata/go-bindata

embed: submodule go-bindata generate build

build:
	go build --ldflags '-s -w'

debug: submodule go-bindata generate
	go build

test:
	go test -v ./...

get-deps:
	go get -t -v ./...

clean:
	rm -f beehive
.PHONY: clean embed go-bindata get-deps noembed generate submodule build all
