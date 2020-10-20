VERSION := $(shell git describe --always --dirty --tags)

.PHONY: default get test all vet

default: get
	CGO_ENABLED=0 go build -i -ldflags="-X 'main.Version=${VERSION}'" -o ${GOPATH}/bin/justget


all: fmt get vet default


get:
	go get

fmt:
	go fmt ./...

dynamic: get
	CGO_ENABLED=1 go build -i -ldflags="-X 'main.Version=${VERSION}'" -o ${GOPATH}/bin/justget

vet:
	go vet ./...

