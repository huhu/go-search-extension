package controller

import (
	"fmt"
	"sync"

	"github.com/huhu/golang/crawler/pkg/data"
	"github.com/huhu/golang/crawler/pkg/gpool"
	"github.com/huhu/golang/crawler/pkg/http"
)

var (
	defaultSince = "2019-04-10T19:08:52.997264Z"
	defaultLimit = 2000
)

type CtlCacheMeta struct {
	Pool    *gpool.Pool
	httpCli *http.CliHTTP
	ChMeta  chan *data.Pkg
	//ConsumeDb *db.SqlitDB
	//ProduceDb *db.SqlitDB
}

func GetCtlCacheMeta() SvcProdCons {
	m := &CtlCacheMeta{}
	m.httpCli = http.NewCliHTTP(true)
	m.ChMeta = make(chan *data.Pkg, 100)

	return m
}

func (m *CtlCacheMeta) GetChan() chan *data.Pkg {
	return m.ChMeta
}

func (m *CtlCacheMeta) Produce() {
	since := defaultSince
	limit := defaultLimit
	first := true
	for {
		metas, err := m.httpCli.GetMetadatasFromCache(since, limit)
		if err != nil {
			panic(err)
		}
		metasLen := len(metas)
		if metasLen == 0 {
			return
		}

		for i := range metas {
			meta := metas[i]
			if !first && meta.Timestamp == since {
				continue
			}
			m.ChMeta <- &data.Pkg{
				FullPath: meta.Path,
				Version:  meta.Version,
			}
		}
		since = metas[metasLen-1].Timestamp
		first = false
		fmt.Printf("since: %s; limit: %d\n", since, limit)
	}
}

func (m *CtlCacheMeta) Consume(wgsql sync.WaitGroup) {

}

func (m *CtlCacheMeta) Start() {

}
