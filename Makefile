GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_TAG = $(shell git describe --tags 2>/dev/null || echo "v.0.0.1" )

PKG        := ./...
LDFLAGS    := -w -s
LDFLAGS += -X "github.com/itxaka/luet-mtree/internal/version.version=${GIT_TAG}"
LDFLAGS += -X "github.com/itxaka/luet-mtree/internal/version.gitCommit=${GIT_COMMIT}"

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

lint: fmt vet


all: lint test build