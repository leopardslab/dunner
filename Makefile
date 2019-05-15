.PHONY: all install test build clean

SHA=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)
TAG?=$(shell git tag -l --contains HEAD)
VERSION=$(TAG)

ifeq ($(VERSION),)
VERSION := latest
endif


#Go parameters
GOCMD=go
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
DEP=dep

all: build

install:
	@$(DEP) ensure

test: install
	@$(GOTEST) -v ./...

build:	install
	@$(GOINSTALL) -ldflags "-X main.version=$(VERSION)-$(SHA) -s"

clean:
	rm -rf *
