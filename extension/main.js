const c = new Compat();

const defaultSuggestion = `Search Go docs in your address bar instantly!`;
const omnibox = new Omnibox(c.browser, defaultSuggestion, c.isChrome ? 8 : 6);

omnibox.bootstrap({
    onSearch: (query) => {
        return [];
    },
    onFormat: (index, doc) => {
    },
    onAppend: () => {
    }
});