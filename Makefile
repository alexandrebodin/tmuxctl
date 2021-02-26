release: 
	@echo "Building releases"
	@goreleaser --rm-dist -f .goreleaser.yml
	@echo "Releases created"
.PHONY: release

clean:
	@echo "Clean dist/"
	@rm -rf dist
.PHONY: clean

test:
	@echo "Running test"
	@go test ./...
.PHONY: test

test.ci:
	@echo "Running test"
	@go test -v ./...
.PHONY: test.ci

build:
	@echo "Building"
	@go build
.PHONY: build

docs:
	@echo "Building docs"
	@static-docs --in docs --out build
.PHONY: docs