package controller

import (
	"fmt"

	"github.com/huhu/golang/crawler/pkg/data"
	"github.com/huhu/golang/crawler/pkg/db"
	"github.com/huhu/golang/crawler/pkg/gpool"
	"github.com/huhu/golang/crawler/pkg/http"

	"log"
	"regexp"
	"runtime"
	"sync"
	"time"
)

type CtlMeta struct {
	Pool      *gpool.Pool
	Cli       *http.CliHTTP
	ChMeta    chan *data.Pkg
	ConsumeDb *db.SqlitDB
	ProduceDb *db.SqlitDB
}

func GetCtlMeta(dbPath string) SvcProdCons {
	m := &CtlMeta{}

	m.ProduceDb = db.GetDB(dbPath).InitDB()
	m.ConsumeDb = db.GetDB(dbPath).InitDB()
	m.ProduceDb.SetRow(data.KeyPathRead, data.SqlPathRead)
	m.ConsumeDb.SetStmt(data.KeyPkgInsert, data.SqlPkgInsert)
	m.Cli = http.NewCliHTTP(false)
	m.ChMeta = make(chan *data.Pkg, 100)

	return m
}

func (m *CtlMeta) GetChan() chan *data.Pkg {
	return m.ChMeta
}

func (m *CtlMeta) Produce() {
	vendorReg := regexp.MustCompile(`(^|.)(vendor).`)
	var id int

	for {
		var str string
		id = id + 1
		nid := m.ProduceDb.GetDataFromDB(data.KeyPathRead, id, &str)
		if nid == -2 {
			break
		}
		if nid == -1 {
			continue
		}
		path := str
		if path == "" {
			break
		}

		m.Pool.Add(1)

		go func(i int, path string, chMeta chan *data.Pkg) {
			// 从 godoc 的 api 获取元数据
			metas := m.Cli.GetMetadatas(path)

			for _, meta := range metas {
				if vendorReg.MatchString(meta.Path) {
					continue
				}
				if meta.Path == "github.com/jbrukh/bayesian" {
					log.Println(meta)
				}

				if meta.Stars >= 500 {
					if meta.Path == path {
						chMeta <- &data.Pkg{
							FullPath:   meta.Path,
							SimpleSpec: meta.Synopsis,
							Stars:      meta.Stars,
							Version:    meta.Version,
						}
						log.Printf("id: %d sent ; num of goroutine:%d\n", i, runtime.NumGoroutine())
					}
				}
			}
			m.Pool.Done()
		}(id, path, m.ChMeta)
	}
}

// Consume: save pkg to database
func (m *CtlMeta) Consume(wgSql sync.WaitGroup) {
	wgSql.Add(1)

	for pkg := range m.ChMeta {
		if pkg.Version == "-c-c-" {
			break
		}

		if id := m.ConsumeDb.SaveMetadata(data.KeyPkgInsert, *pkg); id == -1 {
			log.Printf("Error: save %s failed, return. \n", pkg.FullPath)
			wgSql.Done()
			return
		}
	}

	wgSql.Done()
}

func (m *CtlMeta) Start() {
	ts := time.Now()
	var wgSql sync.WaitGroup
	go m.Consume(wgSql)
	m.Pool = gpool.New(100)
	log.Printf("num of goroutine:%d\n", runtime.NumGoroutine())
	m.Produce()
	time.Sleep(2 * time.Second)
	m.Pool.Wait()
	m.ChMeta <- &data.Pkg{Version: "-c-c-"}
	wgSql.Wait()
	fmt.Println(time.Since(ts))
}
