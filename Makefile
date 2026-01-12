PROJECT_NAME ?= backend-example
MAIN ?= .
BIN_DIR ?= bin
GO ?= go
LDFLAGS ?=
ARGS ?=

.DEFAULT_GOAL := build

.PHONY: all build run test cover clean

all: build

build:
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 $(GO) build -trimpath -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(PROJECT_NAME) $(MAIN)

run:
	$(GO) run $(MAIN) $(ARGS)

test:
	GOTOOLCHAIN=go1.25.5+auto $(GO) test -cover ./...

clean:
	$(RM) -r $(BIN_DIR)