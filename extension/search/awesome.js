function AwesomeSearch(index) {
    this.awesomeIndex = {};
    this.awesomeNames = [];
    index.forEach(([name, url, description, category]) => {
        name = category + ';' + name;
        this.awesomeIndex[name] = {url, description};
        this.awesomeNames.push(name);
    });
}

AwesomeSearch.prototype.search = function(keyword) {
    let result = [];
    keyword = keyword.replace(/[-_$]/g, "");
    for (let rawName of this.awesomeNames) {
        let [category, name] = rawName.split(';');

        let categoryMatchIndex = category.toLowerCase().replace(/[-_$]/g, "").indexOf(keyword);
        if (categoryMatchIndex === -1) {
            categoryMatchIndex = 999;
        }
        let nameMatchIndex = name.toLowerCase().replace(/[-_$]/g, "").indexOf(keyword);
        if (nameMatchIndex === -1) {
            nameMatchIndex = 999;
        }

        if ([categoryMatchIndex, nameMatchIndex].some(i => i !== 999)) {
            result.push({
                name, category, nameMatchIndex, categoryMatchIndex,
                ...this.awesomeIndex[rawName],
            });
        }
    }

    return result.sort((a, b) => {
        if (a.categoryMatchIndex === b.categoryMatchIndex) {
            if (a.nameMatchIndex === b.nameMatchIndex) {
                return a.name.length - b.name.length;
            }
            return a.nameMatchIndex - b.nameMatchIndex;
        }
        return a.categoryMatchIndex - b.categoryMatchIndex;
    });
}