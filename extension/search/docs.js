function DocSearch(searchIndex) {
    this.docs = [];
    Object.entries(searchIndex).forEach(([pkg, docs]) => {
        this.docs = this.docs.concat(docs.map(([label, description, type]) => {
            return {
                description, type, label,
                package: pkg,
            }
        }));
    });
}

DocSearch.prototype.search = function(query) {
    query = query.replace(/[-_\s]/ig, "").toLowerCase();
    let results = [];
    for (let doc of this.docs) {
        let packageMatchIndex = doc.package.indexOf(query);
        if (packageMatchIndex === -1) {
            packageMatchIndex = 999;
        }

        let labelMatchIndex = doc.label.toLowerCase().indexOf(query);
        if (labelMatchIndex === -1) {
            labelMatchIndex = 999;
        }

        if ([packageMatchIndex, labelMatchIndex].some(i => i !== 999)) {
            results.push({
                packageMatchIndex,
                labelMatchIndex,
                ...doc,
            })
        }
    }
    return results.sort((a, b) => {
        if (a.packageMatchIndex === b.packageMatchIndex) {
            if (a.labelMatchIndex === b.labelMatchIndex) {
                return a.label.length - b.label.length;
            }
            return a.labelMatchIndex - b.labelMatchIndex;
        }
        return a.packageMatchIndex - b.packageMatchIndex;
    });
};