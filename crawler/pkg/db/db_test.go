package db

import (
	"log"
	"testing"

	"github.com/huhu/golang/crawler/pkg/data"
)

func TestSaveMetadata(t *testing.T) {

	dbPath := "../svc/test.db"

	t.Run("SaveMetadata", func(t *testing.T) {

		GetDB(dbPath).InitDB()
		db := GetDB(dbPath)
		db.SetStmt(data.KeyPkgInsert, data.SqlPkgInsert)
		p := data.Pkg{
			Libtype:    "name1",
			FullPath:   "path1",
			SimpleSpec: "spec1",
			Stars:      1,
			Version:    "v1",
		}
		num := db.SaveMetadata(data.KeyPkgInsert, p)
		log.Println(num)
	})

	t.Run("InitDB", func(t *testing.T) {
		dbPath := "./test.db"

		GetDB(dbPath).InitDB()
	})

	t.Run("QueryFromDB", func(t *testing.T) {
		db := GetDB(dbPath).InitDB()
		db.SetRow(data.KeyPathRead, data.SqlPathRead)

		var str string
		id := db.GetDataFromDB(data.KeyPathRead, 0, &str)
		log.Printf("id: %d; str: %s\n", id, str)
	})

	t.Run("RunQuerySql", func(t *testing.T) {
		db := GetDB(dbPath).InitDB()
		var str string

		rows, err := db.Query("SELECT * FROM pathes WHERE uid=10")
		if err != nil {
			log.Printf("set query sql failed cuz - %s", err)
			t.Fatal(err)
		}
		if rows.Next() {
			var nuid int
			err := rows.Scan(&nuid, &str)
			if err != nil {
				panic(err)
			}
			log.Printf("nuid: %d; str: %s\n", nuid, str)
		}
	})
}
