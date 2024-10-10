GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
GOMOD = $(GOCMD) mod
GOCYCLO = $(GOPATH)/bin/gocyclo
GOCOGNIT = $(GOPATH)/bin/gocognit
DSN := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

#стандартный ввод т.е просто "make"
all: lint build

build-race:
	$(GOBUILD) --race -o app cmd/app/main.go

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

run-race:
	$(GORUN) cmd/app/main.go --race

run:
	$(GORUN) cmd/app/main.go

deps:
	$(GOMOD) tidy

lint:
	golangci-lint run || true
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
	docker-compose up -d postgres-master

compose-down:
	docker-compose down

compose-ps:
	docker-compose ps

compose-start:
	docker-compose start

compose-stop:public
	docker-compose stop



goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

goose-add:
	goose -dir ./migrations postgres $(DSN) create rename sql

goose-down:
	goose -dir ./migrations postgres $(DSN) reset

goose-up:
	goose -dir ./migrations postgres $(DSN) up

goose-status:
	goose -dir ./migrations postgres $(DSN) status


.PHONY:all build deps run  build-linux build-mac build-windows сlean lint cleanstorages coverage coverage-html coverage-cobertura compose-up compose-down compose-ps compose-start compose-stop
.PHONY: goose-install goose-add goose-up goose-status goose-dowm  compose-up-test compose-down-test run-race build-race