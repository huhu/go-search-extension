package controller

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/huhu/golang/crawler/pkg/crawler"
	"github.com/huhu/golang/crawler/pkg/data"
	"github.com/huhu/golang/crawler/pkg/db"
	"github.com/huhu/golang/crawler/pkg/gpool"
)

var (
	sid = 31513
	eid = 192437
)

type CtlUpdate struct {
	Pool      *gpool.Pool
	ConsumeDb *db.SqlitDB
	ProduceDb *db.SqlitDB
	ChMeta    chan *data.Pkg
}

func GetCtlUpdate(dbPath string) SvcProdCons {
	v := &CtlUpdate{}

	v.ProduceDb = db.GetDB(dbPath).InitDB()
	v.ConsumeDb = db.GetDB(dbPath).InitDB()
	v.ProduceDb.SetRow(data.KeyPkgRead, data.SqlPkgRead)
	v.ProduceDb.SetRow(data.KeyHotPkgRead, data.SqlHotPkgRead)
	v.ConsumeDb.SetStmt(data.KeyPkgUpdate, data.SqlPkgUpdate)

	v.ChMeta = make(chan *data.Pkg, 1000)

	return v
}

func (m *CtlUpdate) GetChan() chan *data.Pkg {
	return m.ChMeta
}

func getMeta(pkg *data.Pkg, flag int) string {
	ch := make(chan string, 1)
	defer close(ch)

	crawler.GetDataFromPkgGo(pkg.FullPath, ch, flag)
	for {
		select {
		case <-time.After(20 * time.Second):
			return ""
		case meta, ok := <-ch:
			if !ok || meta == "-e-r-r-" {
				fmt.Errorf("error meta")
				break
			}

			return meta
		}
	}
	return ""
}

func (v *CtlUpdate) Produce() {
	id := sid
	for {
		hotpkg := &data.Pkg{}
		id = id + 1

		nid := v.ProduceDb.GetDataFromDB(data.KeyPkgRead, id, hotpkg)
		if nid == -2 {
			break
		}
		if nid == -1 {
			break
		}
		if id > eid {
			break
		}

		if hotpkg.Libtype != "" && hotpkg.SimpleSpec != "" && hotpkg.Version != "" {
			continue
		}

		colly := func(i int, pkg *data.Pkg, chMeta chan *data.Pkg) {
			defer v.Pool.Done()
			var j int
			for j = 0; j < 3; j++ {
				if pkg.Version == "" {
					pkg.Version = getMeta(pkg, 1)
				}
				if pkg.Libtype == "" {
					pkg.Libtype = getMeta(pkg, 2)
				}
				if pkg.SimpleSpec == "" {
					pkg.SimpleSpec = getMeta(pkg, 3)
				}

				if pkg.Libtype == "" || pkg.Version == "" {
					log.Printf("sleep 3s and try %d time caz %s libtype or version == \"\": %s ; %s .\n", j,
						pkg.FullPath,
						pkg.Libtype,
						pkg.Version)
					time.Sleep(3 * time.Second)
					continue
				} else {
					break
				}
			}
			if j == 3 {
				return
			}

			chMeta <- pkg
		}

		v.Pool.Add(1)

		go colly(id, hotpkg, v.ChMeta)
	}
	v.Pool.Wait()
}

func (v *CtlUpdate) Consume(wgsql sync.WaitGroup) {
	wgsql.Add(1)
	for meta := range v.ChMeta {
		if meta.Version == "-c-c-" {
			log.Println("break")
			break
		}
		v.ConsumeDb.UpdateMetadata(data.KeyPkgUpdate, meta)
		log.Printf("update data: %s %s  %s %s", meta.Libtype, meta.Version, meta.FullPath, meta.SimpleSpec)

	}

	wgsql.Done()
}

func (v *CtlUpdate) Start() {
	var wgSql sync.WaitGroup
	go v.Consume(wgSql)
	go v.Consume(wgSql)
	go v.Consume(wgSql)

	v.Pool = gpool.New(5000)
	v.Produce()
	v.ChMeta <- &data.Pkg{Version: "-c-c-"}
	wgSql.Wait()
}
