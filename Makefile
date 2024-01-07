# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

# Build target
build:
	$(GOBUILD) -o rss-wmb

# Clean target
clean:
	$(GOCLEAN)
	rm -f rss-wmb

# Test target
test:
	$(GOTEST) -v ./...

# Get dependencies target
deps:
	$(GOGET) -v ./...

.PHONY: build clean test deps
