GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
GOMOD = $(GOCMD) mod
GOCYCLO = $(GOPATH)/bin/gocyclo
GOCOGNIT = $(GOPATH)/bin/gocognit

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

all: lint build

build:
	$(GOBUILD) -o app cmd/app/main.go

build-linux:
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o appL cmd/app/main.go
build-mac:
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o appM cmd/app/main.go
build-windows:
	GOOS=windows GOARCH=arm64 $(GOBUILD) -o appW.exe cmd/app/main.go

clean:
	rm -f ./app ./appL ./appM ./appW.exe

run:
	$(GORUN) cmd/app/main.go

deps:
	$(GOMOD) tidy

#true чтобы прога сбилдилась
lint:
	$(GOCYCLO) -over 5 . || true
	$(GOCOGNIT) -over 5 . || true



.PHONY:all build deps run  build-linux build-mac build-windows сlean lint
