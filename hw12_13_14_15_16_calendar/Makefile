BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

#Postrges
POSTGRES_USER ?= postgres
POSTGRES_PASSWORD ?= password
POSTGRES_DB ?= backend
POSTGRES_PORT ?= 5435
POSTGRES_CONTAINER := postgres-calendar

#Rabbit
RABBIT_USER ?= rabbit
RABBIT_PASSWORD ?= password
RABBIT_PORT ?= 5672
RABBIT_PORT_UI ?= 15672
RABBIT_CONTAINER := rabbitmq


#Protoc
PROTOC_VERSION=36.1

TEST_DATABASE_DSN="postgres://user:userpasswd@localhost:5432/calendar?sslmode=disable&search_path=calendar"


run-postgres:
	docker run -d --name $(POSTGRES_CONTAINER) \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-e POSTGRES_DB=$(POSTGRES_DB) \
	-p $(POSTGRES_PORT):5432 \
	-v postgres-data:/var/lib/postgresql/data \
	postgres:latest

run-rabbit:
	docker run -d --name $(RABBIT_CONTAINER) \
	-p $(RABBIT_PORT):5672 \
	-p $(RABBIT_PORT_UI):15672 \
	-e RABBITMQ_DEFAULT_USER=$(RABBIT_USER) \
	-e RABBITMQ_DEFAULT_PASS=$(RABBIT_PASSWORD) \
	rabbitmq:management

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

install-fmt-deps:
	(which gofumpt > /dev/null) || go install mvdan.cc/gofumpt@latest

fmt: install-fmt-deps
	gofmt -w ./
	gofumpt -l -w ./

run: build
	$(BIN) --config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test: start_db run_tests stop_db
	

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.13.2

lint: install-lint-deps
	golangci-lint run ./...


start_db:
	docker compose -f ./deployments/docker-compose-test.yaml run --rm schema-init

run_tests:
	export TEST_DATABASE_DSN=$(TEST_DATABASE_DSN); go test -race ./internal/... ./cmd/...

stop_db:
	docker compose -f deployments/docker-compose-test.yaml down -v

protoc-deps:
	(which protoc > /dev/null) || curl -fsSL "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip" -o /tmp/protoc.zip && unzip -oq -d $(shell go env GOPATH) /tmp/protoc.zip

protoc: protoc-deps
	protoc --go-grpc_out=./internal/server/rpc/rpcapi --go_out=./internal/server/rpc/rpcapi ./internal/server/rpc/protobuf/calendar.proto



.PHONY: build run build-img run-img version test lint
