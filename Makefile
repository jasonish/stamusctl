PACKAGE=git.stamus-networks.com/lanath/stamus-ctl/internal/app
LOGGER=git.stamus-networks.com/lanath/stamus-ctl/internal/logging

CURRENT_DIR=$(shell pwd)
DIST_DIR=${CURRENT_DIR}/dist
CLI_NAME=stamusctl
DAEMON_NAME=stamusd


HOST_OS:=$(shell go env GOOS)
HOST_ARCH:=$(shell go env GOARCH)

TARGET_ARCH?=linux/amd64

VERSION=$(shell cat ${CURRENT_DIR}/VERSION)
BUILD_DATE:=$(if $(BUILD_DATE),$(BUILD_DATE),$(shell date -u +'%Y-%m-%dT%H:%M:%SZ'))
GIT_COMMIT:=$(if $(GIT_COMMIT),$(GIT_COMMIT),$(shell git rev-parse HEAD))
GIT_TAG:=$(if $(GIT_TAG),$(GIT_TAG),$(shell if [ -z "`git status --porcelain`" ]; then git describe --exact-match --tags HEAD 2>/dev/null; fi))

GOPATH?=$(shell if test -x `which go`; then go env GOPATH; else echo "$(HOME)/go"; fi)
GOCACHE?=$(HOME)/.cache/go-build


STATIC_BUILD?=true

DEV_IMAGE?=false

override LDFLAGS += \
  -X ${PACKAGE}.Arch=${TARGET_ARCH} \
  -X ${PACKAGE}.Commit=${GIT_COMMIT} \
  -X ${PACKAGE}.Version=${VERSION} \
  -X ${LOGGER}.envType=prd

all: cli daemon

cli:
	CGO_ENABLED=0 GODEBUG="tarinsecurepath=0,zipinsecurepath=0" go build -v -ldflags '${LDFLAGS}' -ldflags="-extldflags=-static" -o ${DIST_DIR}/${CLI_NAME} ./cmd

test-cli:
	CGO_ENABLED=0 GODEBUG="tarinsecurepath=0,zipinsecurepath=0" BUILD_MODE=test STAMUS_APP_NAME=stamusctl go build -v -ldflags '${LDFLAGS}' -ldflags="-extldflags=-static" -o ${DIST_DIR}/${CLI_NAME} ./cmd

daemon:
	CGO_ENABLED=0 GODEBUG="tarinsecurepath=0,zipinsecurepath=0" go build -v -ldflags '${LDFLAGS}' -ldflags="-extldflags=-static" -o ${DIST_DIR}/${DAEMON_NAME} ./cmd

daemon-dev:
	air run

test:
	go test ./...

deamon-test:
	go test ./.test/unit

.PHONY: all cli daemon test
