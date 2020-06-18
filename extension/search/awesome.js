function AwesomeSearch(index) {
    this.awesomeIndex = {};
    this.awesomeNames = [];
    index.forEach(([name, url, description, category]) => {
        this.awesomeIndex[name] = {url, description, category};
        this.awesomeNames.push(name);
    });
}

AwesomeSearch.prototype.search = function(keyword) {
    let result = [];
    keyword = keyword.replace(/[-_$]/g, "");
    for (let rawName of this.awesomeNames) {
        let name = rawName.replace(/[-_$]/g, "");
        if (name.length < keyword.length) continue;

        let index = name.indexOf(keyword);
        if (index > -1) {
            result.push({name: rawName, matchIndex: index});
        }
    }

    return result.sort((a, b) => {
        if (a.matchIndex === b.matchIndex) {
            return a.name.length - b.name.length;
        }
        return a.matchIndex - b.matchIndex;
    }).map(item => {
        return {
            name: item.name,
            ...this.awesomeIndex[item.name]
        }
    });
}