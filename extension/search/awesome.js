function AwesomeSearch(index) {
    this.awesomeIndex = {};
    index.forEach(([name, url, description, category]) => {
        let key = category + ';' + name;
        if (!this.awesomeIndex.hasOwnProperty(key)) {
            this.awesomeIndex[key] = [];
        }
        this.awesomeIndex[key].push({category, name, url, description});
    });
    this.awesomeNames = Object.keys(this.awesomeIndex);
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
                rawName, name, nameMatchIndex, categoryMatchIndex,
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
    }).flatMap(item => {
        return this.awesomeIndex[item.rawName];
    });
}