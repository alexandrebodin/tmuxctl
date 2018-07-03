release: 
	@echo "Building releases"
	@goreleaser --rm-dist -f .goreleaser.yml
	@echo "Releases created"
.PHONY: release

clean:
	@echo "Clean dist/"
	@rm -rf dist
.PHONY: clean