package controller

import (
	"sync"
	"time"

	"github.com/huhu/golang/crawler/pkg/data"
	"github.com/huhu/golang/crawler/pkg/http"
)

type SvcProdCons interface {
	Produce()
	Consume(sync.WaitGroup)
	Start()
	GetChan() chan *data.Pkg
}

type CtlPipe struct {
	ctlCache  *CtlCacheMeta
	ctlMeta   SvcProdCons
	ctlUpdate SvcProdCons
}

func (p *CtlPipe) CallUpdate(wg sync.WaitGroup) {
	wg.Add(1)
	time.Sleep(3 * time.Second)
	p.ctlUpdate.Start()
	wg.Done()
}

func (p *CtlPipe) Produce() {
	p.ctlCache.Produce()
}

func (p *CtlPipe) Consume(wg sync.WaitGroup) {
	p.ctlMeta.Consume(wg)
}

// Start Get data from official cache;
// save data to db;
// crawl from go.dev and  update data in db;
func (p *CtlPipe) Start() {
	var wgSql sync.WaitGroup
	go p.Consume(wgSql)
	go p.CallUpdate(wgSql)

	p.Produce()

	time.Sleep(2 * time.Second)

	p.ctlCache.ChMeta <- &data.Pkg{Version: "-c-c-"}
	wgSql.Wait()
}

func (p *CtlPipe) GetChan() chan *data.Pkg {
	return p.ctlMeta.GetChan()
}

func Pipeline(dbPath string) *CtlPipe {

	p := &CtlPipe{}

	p.ctlCache = &CtlCacheMeta{}
	p.ctlCache.httpCli = http.NewCliHTTP(true)
	p.ctlMeta = GetCtlMeta(dbPath)
	p.ctlCache.ChMeta = p.ctlMeta.GetChan()
	p.ctlUpdate = GetCtlUpdate(dbPath)

	return p
}
