package controller

import (
	"sync"
	"testing"

	"github.com/huhu/golang/crawler/pkg/http"
)

func TestPipeline(t *testing.T) {
	dbPath := "./test.db"

	t.Run("CtlPipeGetAndSaveMeta", func(t *testing.T) {
		p := &CtlPipe{}

		p.ctlCache = &CtlCacheMeta{}
		p.ctlCache.httpCli = http.NewCliHTTP(true)
		p.ctlMeta = GetCtlMeta(dbPath)
		p.ctlCache.ChMeta = p.ctlMeta.GetChan()

		p.Start()
	})

	t.Run("Pipeline", func(t *testing.T) {
		p := Pipeline(dbPath)
		p.Start()
	})

	t.Run("Pipeline2", func(t *testing.T) {
		p := Pipeline(dbPath)

		wgSql := sync.WaitGroup{}

		p.CallUpdate(wgSql)
		wgSql.Wait()
	})

}
