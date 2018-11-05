BEEHIVE_VERSION=0.2.99
COMMIT_SHA=$(shell git rev-parse --short HEAD)

all: embed

submodule:
	@echo "Downloading and updating Beehive's submodules..."
	@echo -n "...executing: "
	git submodule update --init
	@echo "Done."

noembed: submodule build
	@echo "Resource embedding: disabled.\n\
	You can use 'make embed' to enable it."

generate:
	@echo "Embedding binary data in go files..."
	@echo -n "...executing: "
	$(shell go env GOPATH)/bin/go-bindata --tags embed --pkg api -o api/bindata.go --ignore \\.git assets/... config/...
	@echo "Data embedded."

go-bindata:
	@echo -n "Management of go-bindata..."
	@if [ ! -f $(shell go env GOPATH)/bin/go-bindata ]; then\
		GOBINDATA_GET="go get -u github.com/jteeuwen/go-bindata/go-bindata";\
		echo "\nDownloading go-bindata for resource embedding...";\
		echo -n "...executing: $$GOBINDATA_GET";\
		$$GOBINDATA_GET;\
		echo "Done.";\
	else\
		echo " nothing to do, go-bindata is already installed.";\
	fi

embed: submodule go-bindata generate build
	@echo "Resource embedding: enabled.\n\
	You can use 'make noembed' to disable it."

build:
	@echo "Building Beehive for production..."
	@echo -n "...executing: "
	go build -tags 'embed' -ldflags '-s -w -X main.Version=$(BEEHIVE_VERSION) -X main.CommitSHA=$(COMMIT_SHA)'
	@echo "Build completed."

debug: submodule go-bindata generate
	@echo "Building Beehive for debugging..."
	@echo -n "...executing: "
	go build -tags 'embed' -ldflags '-X main.Version=$(BEEHIVE_VERSION) -X main.CommitSHA=$(COMMIT_SHA)'
	@echo "Build completed."

test:
	go test -v $(shell go list ./... | grep -v vendor/)

get-deps:
	@echo "Installing Beehive's dependencies..."
	@echo -n "...executing: "
	dep ensure
	@echo "Done."

clean:
	@echo "Removing the Beehive's binary..."
	@echo -n "...executing: "
	rm -f beehive
	@echo "Done. Ready for a new compilation."

.PHONY: clean embed go-bindata get-deps noembed generate submodule build all
