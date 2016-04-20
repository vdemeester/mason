.PHONY: all

MASON_ENVS := \
	-e DOCKER_TEST_HOST \
	-e TESTFLAGS

BIND_DIR := "dist"
MASON_MOUNT := -v "$(CURDIR)/$(BIND_DIR):/go/src/github.com/vdemeester/mason/$(BIND_DIR)"

GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
MASON_DEV_IMAGE := mason-dev$(if $(GIT_BRANCH),:$(GIT_BRANCH))
REPONAME := $(shell echo $(REPO) | tr '[:upper:]' '[:lower:]')

DAEMON_VERSION := $(if $(DAEMON_VERSION),$(DAEMON_VERSION),"default")
DOCKER_RUN_MASON := docker run --rm --privileged -it -e DAEMON_VERSION="$(DAEMON_VERSION)" $(MASON_ENVS) $(MASON_MOUNT) "$(MASON_DEV_IMAGE)"

default: all

all: build ## validate all checks, run all tests
	$(DOCKER_RUN_MASON) ./hack/make.sh

binary: build ## compile the mason binary
	$(DOCKER_RUN_MASON) ./hack/make.sh binary

test: build ## run the unit and integration tests
	$(DOCKER_RUN_MASON) ./hack/make.sh test-unit test-integration

test-integration: build ## run the integration tests
	$(DOCKER_RUN_MASON) ./hack/make.sh test-integration

test-unit: build ## run the unit tests
	$(DOCKER_RUN_MASON) ./hack/make.sh test-unit

validate: build ## validate gofmt, golint and go vet
	$(DOCKER_RUN_MASON) ./hack/make.sh validate-gofmt validate-govet validate-golint

lint:
	./hack/make.sh validate-golint

fmt:
	./hack/make.sh validate-gofmt

build: dist
	docker build -t "$(MASON_DEV_IMAGE)" .

shell: build ## start a shell inside the build env
	$(DOCKER_RUN_MASON) /bin/bash

dist:
	mkdir dist

help: ## this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
