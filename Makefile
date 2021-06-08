GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_TAG = $(shell git describe --tags 2>/dev/null || echo "v.0.0.1" )

PKG        := ./...
LDFLAGS    := -w -s
LDFLAGS += -X "github.com/itxaka/luet-mtree/internal/version.version=${GIT_TAG}"
LDFLAGS += -X "github.com/itxaka/luet-mtree/internal/version.gitCommit=${GIT_COMMIT}"

LUET?=/usr/bin/luet
BACKEND?=docker
CONCURRENCY?=1
export ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
COMPRESSION?=gzip
TREE?=./luet-packages
export LUET_BIN?=$(LUET)


build:
	go build -ldflags '$(LDFLAGS)' -o bin/

vet:
	go vet ${PKG}

fmt:
	go fmt ${PKG}

test:
ifneq ($(shell id -u), 0)
	@echo "Tests need to run under root user to download and unpack docker images."
	@exit 1
else
	go test ${PKG} -race -coverprofile=coverage.txt -covermode=atomic
endif

clean-repo:
	rm -rf build/ *.tar *.metadata.yaml

build-repo: clean-repo
	mkdir -p $(ROOT_DIR)/build
	$(LUET) build --all --tree=$(TREE) --destination $(ROOT_DIR)/build --backend $(BACKEND) --concurrency $(CONCURRENCY) --compression $(COMPRESSION)

create-repo:
	$(LUET) create-repo --tree "$(TREE)" \
    --output $(ROOT_DIR)/build \
    --packages $(ROOT_DIR)/build \
    --name "luet-mtree" \
    --descr "Luet mtree official repository" \
    --urls "http://localhost:8000" \
    --tree-compression gzip \
    --type http


lint: fmt vet


all: lint test build