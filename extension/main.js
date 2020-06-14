const c = new Compat();
const searcher = new DocSearch(searchIndex);
const pkgSearcher = new PackageSearch(pkgs);
const commandManager = new CommandManager(
    new HelpCommand(),
    new HistoryCommand(),
);

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
    },
    afterNavigated: (query, result) => {
        HistoryCommand.record(query, result);
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
            content: `https://pkg.go.dev/${join([pkg.domain, pkg.repository, pkg.name])}`,
            description: `Package: ${pkg.domain}/${c.match(join([pkg.repository, pkg.name]))} ${pkg.version} - ${c.dim(c.escape(pkg.description))}`,
        }
    }
});

// join(["A","bb"]) == "A/bb"
// join(["A","bb",undefined]) == "A/bb"
// join(["A",undefined,undefined]) == "A"
// join(["A",undefined,"a"]) == "A/a"
function join(list) {
    // Use filter() method to filter out falsy item.
    let result = (list || []).filter(_ => _).join("/");
    if (result.endsWith("/")) {
        result = result.slice(0, result.length);
    }
    return result;
}

omnibox.addPrefixQueryEvent(":", {
    onSearch: (query) => {
        return commandManager.execute(query);
    }
});

omnibox.addNoCacheQueries(":");