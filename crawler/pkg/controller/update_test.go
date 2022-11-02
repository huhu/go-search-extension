package controller

import (
	"log"
	"testing"

	"github.com/huhu/golang/crawler/pkg/data"
)

func TestGetCtlVersion(t *testing.T) {
	dbPath := "./test.db"

	t.Run("GetCtlUpdate", func(t *testing.T) {
		ctl := GetCtlUpdate(dbPath)
		ctl.Start()
	})

	t.Run("getMeta", func(t *testing.T) {
		for i := 1; i < 5; i++ {
			meta := getMeta(&data.Pkg{
				FullPath: "github.com/hortonworks/cb-cli",
			}, i)
			log.Printf("meta: %s\n", meta)
		}
	})

}
