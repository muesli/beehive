all: embed

submodule:
	git submodule update --init

noembed: submodule build

generate:
	go-bindata --pkg api -o api/bindata.go --ignore \\.git assets/... config/...

go-bindata:
	[ -f $(shell go env GOPATH)/bin/go-bindata ] || go get -u github.com/jteeuwen/go-bindata/go-bindata

embed: submodule go-bindata generate build

build:
	go build

get-deps:
	go get -u

clean:
	rm -f beehive
.PHONY: clean embed go-bindata get-deps noembed generate submodule build all
