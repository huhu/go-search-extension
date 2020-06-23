function DocSearch(searchIndex) {
    this.docs = searchIndex;
}

DocSearch.prototype.search = function(query) {
    query = query.replace(/[-_\s]/ig, "").toLowerCase();
    let results = [];
    for (let [name, description, type] of this.docs) {
        let matchIndex = name.toLowerCase().indexOf(query);
        if (matchIndex !== -1) {
            results.push({
                name,
                matchIndex,
                description,
                type
            });
        }
    }
    return results.sort((a, b) => {
        if (a.matchIndex === b.matchIndex) {
            return a.name.length - b.name.length;
        }
        return a.matchIndex - b.matchIndex;
    });
};