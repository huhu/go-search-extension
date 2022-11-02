package svc

import (
	"log"
	"regexp"

	"github.com/huhu/golang/crawler/pkg/data"
	"github.com/huhu/golang/crawler/pkg/db"
	"github.com/huhu/golang/crawler/pkg/file"
)

func GetAndSaveMetasFromFile(dbPath, filePath string) {
	d := db.GetDB(dbPath).InitDB()
	defer d.Close()

	d.SetStmt(data.KeySavePath, data.SqlSavePath)

	result := file.GetDataFromFile(filePath)
	vendorReg := regexp.MustCompile(`(^|.)(vendor).`)
	slashReg := regexp.MustCompile(`^[a-zA-Z.0-9\-]+/[a-zA-Z.0-9\-]+/[a-zA-Z.0-9\-]+$`)

	for _, result := range result.Meta {
		if vendorReg.MatchString(result.Path) {
			continue
		}
		if slashReg.MatchString(result.Path) {
			log.Printf("match: %s \n", result.Path)
			d.SaveMetadata(data.KeySavePath, result.Path)
		} else {
		}
		// d.savePkgPath(result.Path)
	}
}
