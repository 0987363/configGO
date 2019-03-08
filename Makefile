PACKAGE = .

export GO15VENDOREXPERIMENT=1

BUILD_VERSION=$(shell git tag|tail -n 1)
BUILD_NUMBER=$(strip $(if $(TRAVIS_BUILD_NUMBER), $(TRAVIS_BUILD_NUMBER), 0))
BUILD_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)

.PHONY: all clean release build run fmt vet lint

SRC = $(shell glide nv $(PACKAGE))

all: run

fmt:
	@rm /tmp/gofmt.log >/dev/null 2>&1 || true
	go fmt $(SRC) ./models | tee -a /tmp/gofmt.log
	@if [ -s /tmp/gofmt.log ]; then false; fi

vet:
	go vet $(SRC) ./models

lint:
	@rm /tmp/golint.log >/dev/null 2>&1 || true
	for dir in $(SRC) ./models; do golint $$dir | tee -a /tmp/golint.log; done
	@if [ -s /tmp/golint.log ]; then false; fi

clean:
	rm -f configGO

release: 
	GO111MODULE=on GOARCH=amd64 GOOS=linux go build -v -ldflags "-X main.BuildVersion=$(BUILD_VERSION).$(BUILD_NUMBER) -X main.BuildCommit=$(BUILD_COMMIT) -X main.BuildDate=$(BUILD_DATE)" 

build: 
	GO111MODULE=on go build -v -ldflags "-X main.BuildVersion=$(BUILD_VERSION).$(BUILD_NUMBER) -X main.BuildCommit=$(BUILD_COMMIT) -X main.BuildDate=$(BUILD_DATE)" 

run: build
	./configGO serve
