# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY=image-builder

all: build
build:
	$(GOBUILD) -o $(BINARY)

test:
	$(GOTEST) ./...

test-verbose:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY)

run: build
	./$(BINARY)

dependency:
	$(GOMOD) download