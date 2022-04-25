NAME ?= demacia
AUTHOR ?= Dinghao
export PACKAGE_NAME ?= $(NAME)
REVISION := $(shell git rev-parse --short=8 HEAD || echo unknown)
BRANCH := $(shell git show-ref | grep "$(REVISION)" | grep -v HEAD | awk '{print $$2}' | sed 's|refs/remotes/origin/||' | sed 's|ref/heads/||' | sort | head -n 1)
BUILT := $(shell date -u +%Y-%m-%dT%H:%M:%S%z)
export TESTFLAGS ?= -cover

LATEST_STABLE_TAG := $(shell git -c versionsort.prereleaseSuffix="-rc" -c versionsort.prereleaseSuffix="-RC" tag -l "v*.*.*" --sort=-v:refname | awk '!/rc/' | head -n 1)
export IS_LATEST :=
ifeq ($(shell git describe --exact-match --match $(LATEST_STABLE_TAG) >/dev/null 2>&1; echo $$?), 0)
export IS_LATEST := true
endif

BUILD_PLATFORMS ?= -os '!netbsd' -os '!openbsd'

PKG = $(PACKAGE_NAME)
COMMON_PACKAGE_NAMESPACE=$(PKG)/core

BUILD_DIR := $(CURDIR)
TARGET_DIR := $(BUILD_DIR)/out

export OUR_PACKAGES ?= $(subst _$(BUILD_DIR),$(PKG),$(shell go list ./... | grep -v '/vendor/'))

ORIGINAL_GOPATH = $(shell echo $$GOPATH | cut -d: -f1)
ORIGINAL_GOPATH_BIN = $(ORIGINAL_GOPATH)/bin

GOPATH_SETUP := $(CURDIR)/.ok

export PATH := $(ORIGINAL_GOPATH_BIN):$(PATH)

export TEST_FLAG_PATH ?= $(subst _$(BUILD_DIR),$(PKG),$(shell find ./../$(PKG) -path ./../app -prune -o -name '*_test.go' -print | sed -e 's|^./../||' | sed -e 's|/[^/]*$./*_test.go||' | uniq))

# Development Tools
LINTER = $(ORIGINAL_GOPATH_BIN)/golangci-lint
DEVELOPMENT_TOOLS = $(LINTER)

define helpCmd
	echo "# Commands: \n\
    	  make all => xxx \n\

endef

env:
	@echo $(ORIGINAL_GOPATH)
	@echo $(GOPATH)
	@echo $(PATH)

help:
	@$(helpCmd)

version:
	@echo Current version: $(VERSION)
	@echo Current revision: $(REVISION)
	@echo Current branch: $(BRANCH)
	@echo Build platforms: $(BUILD_PLATFORMS)
	@echo IS_LATEST: $(IS_LATEST)

gomod_check: $(GOPATH_SETUP)
	go mod verify

gomod_download: $(GOPATH_SETUP)
	go mod download

golint_check: $(GOPATH_SETUP)
	golint ./...

govet_check: $(GOPATH_SETUP)
	go vet ./...

# local GOPATH
$(GOPATH_SETUP):
	touch $@

$(LINTER): $(GOPATH_SETUP)
	go get github.com/golangci/golangci-lint/cmd/golangci-lint

pull_images_for_tests: $(GOPATH_SETUP)
	# Pulling images required for some tests
	@go run ./.gitlab/scripts/pull-images-for-tests/main.go

check_race_conditions:
	@./.gitlab/scripts/check_race_conditions $(TEST_FLAG_PATH)

parallel_test_prepare: $(GOPATH_SETUP)
	# Preparing test commands
	@./.gitlab/scripts/go_test_with_coverage_report prepare

parallel_test_execute: $(GOPATH_SETUP)
	# executing tests
	@./.gitlab/scripts/go_test_with_coverage_report execute

parallel_test_coverage_report: $(GOPATH_SETUP)
	# Preparing coverage report
	@./.gitlab/scripts/go_test_with_coverage_report coverage

parallel_test_junit_report: $(GOPATH_SETUP)
	# Preparing jUnit test report
	@./.gitlab/scripts/go_test_with_coverage_report junit

markdown_report: $(GOPATH_SETUP)
	@./.gitlab/scripts/markdown_report

check_modules:
	# check go.sum
	@git diff go.sum > /tmp/gosum-$${CI_JOB_ID}-before
	@go mod tidy
	@git diff go.sum > /tmp/gosum-$${CI_JOB_ID}-after
	@diff -U0 /tmp/gosum-$${CI_JOB_ID}-before /tmp/gosum-$${CI_JOB_ID}-after
