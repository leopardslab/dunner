#Go parameters

GOCMD=go
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
DEP=dep 
.PHONY : all dep install test

all : def

def :
	@$(GOINSTALL) -ldflags '-s'

install: dep

dep:
	@$(DEP) ensure

test:
	@$(GOTEST) -v ./...

clean:
	rm -rf *