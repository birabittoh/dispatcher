PROJECT_NAME ?= webhook-dispatcher
MAIN ?= .
DIST_DIR ?= dist
GO ?= go
LDFLAGS ?=
ARGS ?=
COVERAGE_FILE ?= coverage.out

.DEFAULT_GOAL := build

.PHONY: all build run test cover clean

all: build

build:
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 $(GO) build -trimpath -ldflags '$(LDFLAGS)' -o $(DIST_DIR)/$(PROJECT_NAME) $(MAIN)

run:
	$(GO) run $(MAIN) $(ARGS)

test:
	$(GO) test -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@echo "Coverage report generated at $(COVERAGE_FILE)"
	@$(GO) tool cover -func=$(COVERAGE_FILE) | grep total

cover: test
	$(GO) tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "Coverage HTML report: coverage.html"

clean:
	$(RM) -r $(DIST_DIR)
	$(RM) $(COVERAGE_FILE) coverage.html
