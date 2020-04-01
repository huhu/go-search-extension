function PackageSearch(packagesIndex) {
    this.packages = packagesIndex;
    this.packageNames = Object.keys(packagesIndex);
}

PackageSearch.prototype.search = function(keyword) {
    let result = [];
    keyword = keyword.replace(/[-_!]/g, "");
    for (let rawName of this.packageNames) {
        let [domain, repository, name] = rawName.replace(/[-_]/ig, "").split("/");
        let domainMatchIndex = domain.indexOf(keyword);
        if (domainMatchIndex === -1) {
            domainMatchIndex = 999;
        }
        let repositoryMatchIndex = repository.indexOf(keyword);
        if (repositoryMatchIndex === -1) {
            repositoryMatchIndex = 999;
        }
        let nameMatchIndex = name.indexOf(keyword);
        if (nameMatchIndex === -1) {
            nameMatchIndex = 999;
        }

        if ([domainMatchIndex, repositoryMatchIndex, nameMatchIndex].some(i => i !== 999)) {
            let [description, version] = this.packages[rawName];
            result.push({
                domainMatchIndex,
                repositoryMatchIndex,
                nameMatchIndex,
                description,
                version,
                repository,
                name,
                domain,
            });
        }
    }
    return result.sort((a, b) => {
        if (a.domainMatchIndex === b.domainMatchIndex) {
            if (a.repositoryMatchIndex === b.repositoryMatchIndex) {
                if (a.nameMatchIndex === b.nameMatchIndex) {
                    return a.name.length - b.name.length;
                }
                return a.nameMatchIndex - b.nameMatchIndex;
            }
            return a.repositoryMatchIndex - b.repositoryMatchIndex;
        }
        return a.domainMatchIndex - b.domainMatchIndex;
    });
};