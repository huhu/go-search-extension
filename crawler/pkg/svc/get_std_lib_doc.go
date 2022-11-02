package svc

import (
	"log"

	"github.com/huhu/golang/crawler/pkg/crawler"
	"github.com/huhu/golang/crawler/pkg/data"
	"github.com/huhu/golang/crawler/pkg/db"
)

func GetAndSaveStdDoc(dbPath, url string) {
	lablesMap := make(map[string]data.Datum)
	pkgsMap := make(map[string]data.Pkg)

	d := db.GetDB(dbPath).InitDB()
	defer d.Close()

	d.SetStmt(data.KeyPkgInsert, data.SqlPkgInsert)
	d.SetStmt(data.KeyDocInsert, data.SqlDocInsert)

	crawler.GetStdLibDoc(url, lablesMap, pkgsMap)

	for _, pkg := range pkgsMap {
		id := d.SaveMetadata(data.KeyPkgInsert, pkg)
		log.Printf("saved datum: %d\n", id)
	}

	for _, doc := range lablesMap {
		id := d.SaveMetadata(data.KeyDocInsert, doc)
		log.Printf("saved standard pkg: %d\n", id)
	}
}
