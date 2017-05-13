all: embed

submodule:
	git submodule update --init

noembed: submodule build

generate:
	go-bindata --tags embed --pkg api -o api/bindata.go --ignore \\.git assets/... config/...

go-bindata:
	[ -f $(shell go env GOPATH)/bin/go-bindata ] || go get -u github.com/jteeuwen/go-bindata/go-bindata

embed: submodule go-bindata generate build

build:
	go build -tags 'embed' -ldflags '-s -w'

debug: submodule test go-bindata generate
	go build -tags 'embed'

test:
	go test -v $(shell go list ./... | grep -v vendor/)

get-deps:
	go get -t -d $(shell go list ./... | grep -v vendor/)

clean:
	rm -f beehive
.PHONY: clean embed go-bindata get-deps noembed generate submodule build all
