#Go parameters
GOCMD=go
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
DEP=dep 

all : dep install

dep:
	@$(DEP) ensure

install:
	@$(GOINSTALL) -ldflags '-s'

test:
	@$(GOTEST) -v ./...