package data

var (
	SqlPkgInsert = `INSERT INTO pkgs(libtype, fullPath,simpleSpec, stars,version) values(?,?,?,?,?)`
	KeyPkgInsert = "pkg_insert"
	SqlDocInsert = `INSERT INTO docs(keyname, label, definition,simpleSpec,datatype,pkg) values(?,?,?,?,?,?)`
	KeyDocInsert = "doc_insert"
	KeySavePath  = "save_metas"
	SqlSavePath  = `INSERT INTO pathes(path) values(?)`
	KeyPathRead  = "path_read"
	SqlPathRead  = "SELECT * FROM pathes WHERE uid=%s"
	KeyPkgRead   = "pkg_read"
	SqlPkgRead   = "SELECT * FROM pkgs WHERE uid=%s"

	KeyHotPkgRead = "hot_read"
	SqlHotPkgRead = "SELECT * FROM hotpkg WHERE uid=%s"

	KeyPkgUpdate = "pkg_update"
	SqlPkgUpdate = `update pkgs set libtype=?,simpleSpec=?,version=? where fullpath=?`

	KeyCreatePkgs = "create_pkgs"
	SqlCreatePkgs = `CREATE TABLE IF NOT EXISTS pkgs
-- the collection of packages
(
uid  INTEGER PRIMARY KEY AUTOINCREMENT,
libtype VARCHAR(64) NULL, -- the type of lib: package / command
fullPath   VARCHAR(64) NULL,	-- the full path of package
simpleSpec  TEXT NULL,	-- the first sentence of the package specification
stars  INTERGER, -- how many stars this package repo got in github
version  VARCHAR(64) NULL -- version of this pkg
);`
	KeyCreatePathes = "create_pathes"
	SqlCreatePathes = `CREATE TABLE IF NOT EXISTS pathes
-- the collection of packages
(
 uid  INTEGER PRIMARY KEY AUTOINCREMENT,
 path   VARCHAR(64) NULL	-- the full path of package
);`
	KeyCreateDocs = "create_docs"
	SqlCreateDocs = `CREATE TABLE IF NOT EXISTS  docs 
-- the collection of datum, a datum is in an unique package
(
 uid  INTEGER PRIMARY KEY AUTOINCREMENT,
 keyname  VARCHAR(64) NULL,	-- the name of datum
 label  VARCHAR(64) NULL,	-- the unique anchor in the web site
 definition  TEXT NULL,		-- the full definition of the datum
 simpleSpec  TEXT NULL,		-- the first sentence of the specification
 datatype  VARCHAR(64) NULL,	-- the datatype of the datum, could be 'interface','struct','func' and 'other'
 pkg  VARCHAR(64) NULL	-- the full path of package where the datum from
);`
)
