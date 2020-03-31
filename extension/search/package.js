function PackageSearch(packagesIndex) {
    this.packages = packagesIndex;
    this.packageNames = Object.keys(packagesIndex);
}

PackageSearch.prototype.search = function(keyword) {
    let result = [];
    keyword = keyword.replace(/[-_!\s]/g, "");
    for (let rawName of this.packageNames) {
        let name = rawName.replace(/[-_\s]/ig, "");
        if (name.length < keyword.length) continue;

        let index = name.indexOf(keyword);
        if (index !== -1) {
            for (let [fullPath, description, version] of this.packages[rawName]) {
                result.push({
                    name: rawName,
                    matchIndex: index,
                    fullPath,
                    description,
                    version,
                });
            }
        }
    }
    return result.sort((a, b) => {
        if (a.matchIndex === b.matchIndex) {
            return a.name.length - b.name.length;
        }
        return a.matchIndex - b.matchIndex;
    });
};