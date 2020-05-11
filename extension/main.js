const c = new Compat();
const searcher = new DocSearch(searchIndex);
const pkgSearcher = new PackageSearch(pkgs);

const defaultSuggestion = `Search Go std docs and third packages in your address bar instantly!`;
const omnibox = new Omnibox(defaultSuggestion, c.omniboxPageSize());

omnibox.bootstrap({
    onSearch: (query) => {
        return searcher.search(query);
    },
    onFormat: (index, doc) => {
        let text = doc.package;
        let path = doc.package;
        if (doc.type !== "package") {
            text += `.${doc.label}`;
            path += `#${doc.label}`;
        }
        return {
            content: `https://pkg.go.dev/${path}`,
            description: `Std docs: [${doc.type}] ${c.match(text)} - ${c.dim(c.escape(doc.description))}`,
        };
    },
    onAppend: (query) => {
        return [{
            content: `https://pkg.go.dev/search?q=${query}`,
            description: `Search packages ${c.match(query)} on https://pkg.go.dev/`,
        }]
    }
});

omnibox.addPrefixQueryEvent("!", {
    defaultSearch: true,
    searchPriority: 1,
    onSearch: (query) => {
        return pkgSearcher.search(query);
    },
    onFormat: (index, pkg) => {
        return {
            content: `https://pkg.go.dev/${pkg.domain}/${pkg.repository}/${pkg.name}`,
            description: `Package: ${pkg.domain}/${c.match(pkg.repository + "/" + pkg.name)} ${pkg.version} - ${c.dim(c.escape(pkg.description))}`,
        }
    }
});