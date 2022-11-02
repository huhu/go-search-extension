package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"

	"github.com/huhu/golang/crawler/pkg/data"
	"github.com/pkg/errors"

	"log"

	_ "github.com/mattn/go-sqlite3"
)

// SqlitDB db object
type SqlitDB struct {
	*sql.DB
	path    string
	stmtMap map[string]func() (*sql.Stmt, error)
	rowMap  map[string]func(int) (*sql.Rows, error)
	sync.Mutex
}
type NullString struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

func (db *SqlitDB) InitDB() *SqlitDB {
	db.SetStmt(data.KeyCreatePathes, data.SqlCreatePathes)
	db.SetStmt(data.KeyCreatePkgs, data.SqlCreatePkgs)
	db.SetStmt(data.KeyCreateDocs, data.SqlCreateDocs)

	if err := db.createTableForSqlite3(data.KeyCreatePathes, nil); err != nil {
		log.Panic(err)
	}

	if err := db.createTableForSqlite3(data.KeyCreatePkgs, nil); err != nil {
		log.Panic(err)
	}

	if err := db.createTableForSqlite3(data.KeyCreateDocs, nil); err != nil {
		log.Panic(err)
	}
	return db
}

func (db *SqlitDB) createTableForSqlite3(key string, itf interface{}) error {
	db.Lock()
	defer db.Unlock()

	switch itf.(type) {
	case nil:
		stmt, err := db.stmtMap[key]()
		if err != nil {
			panic(err)
		}
		res, err := stmt.Exec()
		if err != nil {
			return err
		}
		defer stmt.Close()

		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		log.Printf("create table: %s %d\n", key, id)
		return nil
	}
	return errors.New("create table failed!")
}

func GetDB(dbPath string) *SqlitDB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	return &SqlitDB{
		DB:      db,
		path:    dbPath,
		stmtMap: make(map[string]func() (*sql.Stmt, error)),
		rowMap:  make(map[string]func(int) (*sql.Rows, error)),
		Mutex:   sync.Mutex{},
	}
}

func (db *SqlitDB) SetStmt(key, sqlStr string) error {
	funSetting := func() (*sql.Stmt, error) {
		stmt, err := db.Prepare(sqlStr)
		if err != nil {
			log.Printf("prepare sql failed cuz - %s", err)
			return nil, err
		}
		return stmt, nil
	}

	db.stmtMap[key] = funSetting
	return nil
}

func (db *SqlitDB) SetRow(key, sqlStr string) error {
	funQuery := func(uid int) (*sql.Rows, error) {
		sid := strconv.FormatInt(int64(uid), 10)
		raws, err := db.Query(fmt.Sprintf(sqlStr, sid))
		if err != nil {
			log.Printf("set query sql failed cuz - %s", err)
			return nil, err
		}
		return raws, nil
	}
	db.rowMap[key] = funQuery
	return nil
}

// return uid : normal
// return -1: query no result
// return -2: case miss
func (db *SqlitDB) GetDataFromDB(key string, uid int, itf interface{}) int {
	db.Lock()
	defer db.Unlock()

	switch itf.(type) {
	case *data.PkgVersions:
		d := itf.(*data.PkgVersions)
		rows, err := db.rowMap[key](uid)
		if err != nil {
			log.Printf("query failed: %s \n", err)
			return -1
		}
		defer rows.Close()
		if rows.Next() {
			var nuid int
			err := rows.Scan(&d.ID, &d.VersionNum, &d.FullPath)
			if err != nil {
				panic(err)
			}
			return nuid
		}
	case *data.Pkg:
		d := itf.(*data.Pkg)
		rows, err := db.rowMap[key](uid)
		if err != nil {
			log.Printf("query failed: %s \n", err)
			return -1
		}
		defer rows.Close()

		if rows.Next() {
			var nuid int
			var in interface{}
			var a sql.NullString
			var b sql.NullString
			var c sql.NullString

			err := rows.Scan(&nuid, &a, &d.FullPath, &b, &in, &c)
			if a.Valid {
				// use s.String
				d.Libtype = a.String
			} else {
				// NULL value
				d.Libtype = ""
			}
			if b.Valid {
				// use s.String
				d.SimpleSpec = b.String
			} else {
				// NULL value
				d.SimpleSpec = ""
			}
			if c.Valid {
				// use s.String
				d.Version = c.String
			} else {
				// NULL value
				d.Version = ""
			}
			if err != nil {
				panic(err)
			}
			return nuid
		}
	case *string:
		path := itf.(*string)
		fn := db.rowMap[key]
		rows, err := fn(uid)
		if err != nil {
			log.Printf("query failed: %s \n", err)
			return -1
		}
		defer rows.Close()

		if rows.Next() {
			var nuid int
			err := rows.Scan(&nuid, path)
			if err != nil {
				panic(err)
			}
			return nuid
		}
	}

	log.Printf("miss case: %d", uid)

	return -2
}

func (db *SqlitDB) SaveMetadata(key string, itf interface{}) int64 {
	db.Lock()
	defer db.Unlock()

	switch itf.(type) {
	case data.Pkg:
		d := itf.(data.Pkg)
		stmt, err := db.stmtMap[key]()
		if err != nil {
			panic(err)
		}
		defer stmt.Close()

		res, err := stmt.Exec(d.Libtype, d.FullPath, d.SimpleSpec, d.Stars, d.Version)
		if err != nil {
			panic(err)
		}

		id, err := res.LastInsertId()
		if err != nil {
			panic(err)
		}
		//log.Printf("saved pkg: %s %d\n", d.Libtype, id)
		return id

	case data.Datum:
		d := itf.(data.Datum)
		stmt, err := db.stmtMap[key]()
		if err != nil {
			panic(err)
		}
		defer stmt.Close()

		res, err := stmt.Exec(d.Name, d.Label, d.Definition, d.SimpleSpec, d.Datatype, d.Pkg)
		if err != nil {
			panic(err)
		}

		id, err := res.LastInsertId()
		if err != nil {
			panic(err)
		}
		log.Printf("saved datum: %s %d\n", d.Name, id)
		return id
	case string:
		path := itf.(string)
		stmt, err := db.stmtMap[key]()
		if err != nil {
			panic(err)
		}
		defer stmt.Close()

		res, err := stmt.Exec(path)
		if err != nil {
			panic(err)
		}

		id, err := res.LastInsertId()
		if err != nil {
			panic(err)
		}
		log.Printf("saved pathes metadata: %s %d\n", path, id)
		return id
	}
	return -1
}

func (db *SqlitDB) UpdateMetadata(key string, d *data.Pkg) {
	db.Lock()
	defer db.Unlock()

	////更新数据
	stmt, err := db.stmtMap[key]()
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(d.Libtype, d.SimpleSpec, d.Version, d.FullPath)
	if err != nil {
		panic(err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		panic(err)
	}

	//fmt.Println(affect)
}
