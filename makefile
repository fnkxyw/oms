GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
GOMOD = $(GOCMD) mod
GOCYCLO = $(GOPATH)/bin/gocyclo
GOCOGNIT = $(GOPATH)/bin/gocognit
DSN := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
LOCAL_BIN := $(CURDIR)/bin
OUT_PATH := $(CURDIR)/pkg
#стандартный ввод т.е просто "make"

# запуск сервера из proto-файла
all: deps bin-deps generate build-server run-server

build-race:
	$(GOBUILD) --race -o app cmd/app/main.go

build-server:
	$(GOBUILD) -o server_app cmd/grpc-pup/pup-service/main.go

build-client:
	$(GOBUILD) -o client_app cmd/app/main.go

#лишним не будет как посчитал
build-linux:
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o appL cmd/app/main.go
build-mac:
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o appM cmd/app/main.go
build-windows:
	GOOS=windows GOARCH=arm64 $(GOBUILD) -o appW.exe cmd/app/main.go

#удалять бинарники
clean:
	rm -f ./client_app ./appL ./appM ./appW.exe coverage.out ./server_app

run-race:
	$(GORUN) cmd/app/main.go --race
#запуск сервера
run-server:
	$(GORUN) cmd/grpc-pup/pup-service/main.go

#запуска клиента к серверу
run-client:
	$(GORUN) cmd/app/main.go

run:
	$(GORUN) cmd/grpc-pup/pup-service/main.go & \
	$(GORUN) cmd/app/main.go & \
	wait

deps:
	$(GOMOD) tidy

lint:
	golangci-lint run || true
	squawk migrations/* || true
	protolint api/pup-service/v1/pup_service.proto

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


# --------------------------
# Поднятие Docker-контейнера
# ---------------------------
compose-up:
	docker-compose up -d

compose-down:
	docker-compose down

compose-ps:
	docker-compose ps

compose-start:
	docker-compose start

compose-stop:
	docker-compose stop

# --------------------------
# ---------------------------


# --------------------------
# Работа с миграциями
# ---------------------------
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

# --------------------------
# ---------------------------


bin-deps: .vendor-proto
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@latest

generate:
	mkdir -p ${OUT_PATH}
	protoc --proto_path api --proto_path vendor.protogen \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out=${OUT_PATH} --go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out=${OUT_PATH} --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway --grpc-gateway_out ${OUT_PATH} --grpc-gateway_opt paths=source_relative \
		--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 --openapiv2_out=${OUT_PATH} \
		--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate --validate_out="lang=go,paths=source_relative:${OUT_PATH}" \
		./api/PuP-service/v1/pup_service.proto
	mv ${OUT_PATH}/PuP-service/v1/pup_service.swagger.json $(CURDIR)/cmd/grpc-pup/pup-service/swagger

.vendor-proto: .vendor-proto/google/protobuf .vendor-proto/google/api .vendor-proto/protoc-gen-openapiv2/options .vendor-proto/validate

.vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/grpc-ecosystem && \
 		cd vendor.protogen/grpc-ecosystem && \
		git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
		git checkout
		mkdir -p vendor.protogen/protoc-gen-openapiv2
		mv vendor.protogen/grpc-ecosystem/protoc-gen-openapiv2/options vendor.protogen/protoc-gen-openapiv2
		rm -rf vendor.protogen/grpc-ecosystem

.vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor.protogen/protobuf &&\
		cd vendor.protogen/protobuf &&\
		git sparse-checkout set --no-cone src/google/protobuf &&\
		git checkout
		mkdir -p vendor.protogen/google
		mv vendor.protogen/protobuf/src/google/protobuf vendor.protogen/google
		rm -rf vendor.protogen/protobuf

.vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor.protogen/googleapis && \
 		cd vendor.protogen/googleapis && \
		git sparse-checkout set --no-cone google/api && \
		git checkout
		mkdir -p  vendor.protogen/google
		mv vendor.protogen/googleapis/google/api vendor.protogen/google
		rm -rf vendor.protogen/googleapis

.vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor.protogen/tmp && \
		cd vendor.protogen/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor.protogen/validate
		mv vendor.protogen/tmp/validate vendor.protogen/
		rm -rf vendor.protogen/tmp


.PHONY:all build deps run  build-linux build-mac build-windows сlean lint cleanstorages coverage coverage-html coverage-cobertura compose-up compose-down compose-ps compose-start compose-stop run-server run-client run
.PHONY: goose-install goose-add goose-up goose-status goose-dowm  compose-up-test compose-down-test run-race build-race generate .vendor-proto .vendor-proto/google/api .vendor-proto/google/protobuf .vendor-proto/protoc-gen-openapiv2/options .vendor-proto/validate