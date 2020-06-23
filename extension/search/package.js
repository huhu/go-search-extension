const LOWEST_MATCH_INDEX = 999;

function PackageSearch(packagesIndex) {
    this.packages = packagesIndex;
    this.packageNames = Object.keys(packagesIndex);
}

PackageSearch.prototype.search = function(keyword) {
    let result = [];
    keyword = keyword.replace(/[-_!]/g, "");
    for (let rawName of this.packageNames) {
        let [domain, repository, ...name] = rawName.split("/");
        // Join the rest to name.
        name = name.join("/");
        let domainMatchIndex = domain.replace(/[-_]/ig, "").indexOf(keyword);
        if (domainMatchIndex === -1) {
            domainMatchIndex = LOWEST_MATCH_INDEX;
        }
        let repositoryMatchIndex = repository ? repository.replace(/[-_]/ig, "").indexOf(keyword) : LOWEST_MATCH_INDEX;
        if (repositoryMatchIndex === -1) {
            repositoryMatchIndex = LOWEST_MATCH_INDEX;
        }

        let nameMatchIndex = name ? name.replace(/[-_]/ig, "").indexOf(keyword) : LOWEST_MATCH_INDEX;
        if (nameMatchIndex === -1) {
            nameMatchIndex = LOWEST_MATCH_INDEX;
        }

        if ([domainMatchIndex, repositoryMatchIndex, nameMatchIndex].some(i => i !== LOWEST_MATCH_INDEX)) {
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
                if (a.nameMatchIndex === b.nameMatchIndex && a.name && b.name) {
                    return a.name.length - b.name.length;
                }
                return a.nameMatchIndex - b.nameMatchIndex;
            }
            return a.repositoryMatchIndex - b.repositoryMatchIndex;
        }
        return a.domainMatchIndex - b.domainMatchIndex;
    });
};