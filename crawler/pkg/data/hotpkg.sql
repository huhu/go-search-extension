CREATE TABLE IF NOT EXISTS  hotpkg
-- the collection of hot package
(
 uid  INTEGER PRIMARY KEY AUTOINCREMENT,
 num  VARCHAR(64) NULL,	-- the number of pkg's versions
 fullPath  VARCHAR(64) NULL	-- the full path of package
);

INSERT INTO hotpkg(num, fullPath) SELECT COUNT(version),fullPath FROM pkgs GROUP BY fullPath ORDER BY count(version) DESC;