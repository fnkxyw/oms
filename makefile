GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
GOMOD = $(GOCMD) mod

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

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



.PHONY: build install run  build-linux build-mac build-windows сlean