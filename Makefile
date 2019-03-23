#Go parameters

GOCMD=go
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
DEP=dep 
.PHONY : all install vet fmt test

all : def

def :
	@$(GOINSTALL) -ldflags '-s'

install: 
	@$(DEP) ensure

test:
	@$(GOTEST) -v ./...

clean:
	rm -rf *

vet:
	@echo "=== go vet ==="
	@go vet ./...

fmt:
	@echo "=== go fmt ==="
	@go fmt ./...
