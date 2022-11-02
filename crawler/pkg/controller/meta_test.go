package controller

import (
	"log"
	"runtime"
	"sync"
	"testing"

	"github.com/huhu/golang/crawler/pkg/data"
	"github.com/huhu/golang/crawler/pkg/db"
	"github.com/huhu/golang/crawler/pkg/http"
)

func TestGetCtlMeta(t *testing.T) {
	metas := [5]data.Metadata{
		{
			Name: "a1",
			Path: "a2",
		},
		{
			Name: "b1",
			Path: "b2",
		},
		{
			Name: "c1",
			Path: "c2",
		},
		{
			Name: "d1",
			Path: "d2",
		},
		{
			Name: "-c-c-",
		},
	}

	dbPath := "../svc/test.db"

	t.Run("GetCtlMeta", func(t *testing.T) {
		ctl := GetCtlMeta(dbPath)
		ctl.Start()
	})

	t.Run("TestComsume", func(t *testing.T) {
		m := &CtlMeta{}

		m.ProduceDb = db.GetDB(dbPath).InitDB()
		m.ConsumeDb = db.GetDB(dbPath).InitDB()
		m.ProduceDb.SetRow(data.KeyPathRead, data.SqlPathRead)
		m.ConsumeDb.SetStmt(data.KeyPkgInsert, data.SqlPkgInsert)
		m.Cli = http.NewCliHTTP(false)
		m.ChMeta = make(chan *data.Pkg, 100)

		var wgSql sync.WaitGroup
		go func() {
			for i, meta := range metas {
				m.ChMeta <- &data.Pkg{
					FullPath:   meta.Path,
					SimpleSpec: meta.Synopsis,
					Stars:      meta.Stars,
					Version:    meta.Version,
				}
				log.Printf("id: %d ; num of goroutine:%d\n", i, runtime.NumGoroutine())
			}
		}()
		m.Consume(wgSql)
		wgSql.Wait()
	})

}
