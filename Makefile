.PHONY: bootstrap lint dependency clean build release all

VERSION := "v1.0.0"
COMMIT  := $(shell git describe --always)
PKGS    := $(shell go list ./... | grep -v /vendor)

default: build

bootstrap:
	@echo "Bootstraping..."
	go get -u golang.org/x/lint/golint
	go get -u github.com/tcnksm/ghr
	go get -u github.com/Masterminds/glide

lint:
	@echo "Source Code Lint..."
	@for i in $(PKGS); do echo $${i}; golint $${i}; done

dependency:
	glide install --strip-vendor

build-linux:
	@echo "Creating Build for Linux..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./releases/$(VERSION)/kubectl-search-Linux-x86_64

build-darwin:
	@echo "Creating Build for macOS..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ./releases/$(VERSION)/kubectl-search-Darwin-x86_64

build-windows:
	@echo "Creating Build for Windows..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ./releases/$(VERSION)/kubectl-search-Windows-x86_64.exe

build: build-linux build-darwin build-windows

clean:
	@echo "Cleanup Releases..."
	rm -rf ./releases/*

release:
	@echo "Creating Releases..."
	go get -u github.com/tcnksm/ghr
	ghr -t ${GITHUB_TOKEN} $(VERSION) releases/$(VERSION)/

all: bootstrap lint dependency clean build
