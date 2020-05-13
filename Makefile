.PHONY: manifest

manifest:
	@jsonnet -J core $@.jsonnet --ext-str browser=$(browser) -o extension/$@.json
	@cp -R core/src extension/core
