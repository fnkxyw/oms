GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
GOMOD = $(GOCMD) mod
GOCYCLO = $(GOPATH)/bin/gocyclo
GOCOGNIT = $(GOPATH)/bin/gocognit

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

#стандартный ввод т.е просто "make"
all: lint build

build:
	$(GOBUILD) -o app cmd/app/main.go

#лишним не будет как посчитал
build-linux:
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o appL cmd/app/main.go
build-mac:
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o appM cmd/app/main.go
build-windows:
	GOOS=windows GOARCH=arm64 $(GOBUILD) -o appW.exe cmd/app/main.go

#удалять бинарники
clean:
	rm -f ./app ./appL ./appM ./appW.exe

run:
	$(GORUN) cmd/app/main.go

deps:
	$(GOMOD) tidy

#true чтобы прога сбилдилась даже если будут файлы   >5
lint:
	golangci-lint run

#очистить хранилища если будет необходимость
cleanstorages:
	rm -f api/*.json



.PHONY:all build deps run  build-linux build-mac build-windows сlean lint cleanstorages