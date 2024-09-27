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
	rm -f ./app ./appL ./appM ./appW.exe coverage.out

run:
	$(GORUN) cmd/app/main.go

deps:
	$(GOMOD) tidy

#true чтобы прога сбилдилась даже если будут файлы   >5
lint:
	golangci-lint run
	squawk migrations/*

#очистить хранилища если будет необходимость
cleanstorages:
	rm -f api/*.json

coverage:
	go test -coverprofile=coverage.out ./...


# Команда для расчета покрытия тестов и открытия HTML отчета
coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

coverage-cobertura:
	go test ./... -coverprofile=coverage.txt -covermode=count
	gocover-cobertura < coverage.txt > coverage.xml

compose-up:
	docker-compose up -d postgres

compose-down:
	docker-compose down postgres

compose-ps:
	docker-compose ps

compose-start:
	docker-compose start postgres

compose-stop:
	docker-compose stop postgres

goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

goose-add:
	goose -dir ./migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" create rename sql

goose-down:
	goose -dir ./migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" down

goose-up:
	goose -dir ./migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up

goose-status:
	goose -dir ./migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" status


.PHONY:all build deps run  build-linux build-mac build-windows сlean lint cleanstorages coverage coverage-html coverage-cobertura compose-up compose-down compose-ps compose-start compose-stop
.PHONY: goose-install goose-add goose-up goose-status goose-dowm