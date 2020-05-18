.PHONY: manifest

manifest: clean
	@jsonnet -J core $@.jsonnet --ext-str browser=$(browser) -o extension/$@.json
	@cp -R core/src extension/core

clean:
	@rm -rf extension/core manifest.json
