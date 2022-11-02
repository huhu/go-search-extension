package crawler

import (
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"strings"
)

var godevURL = "https://pkg.go.dev/"


func GetDataFromPkgGo(path string, ch chan string, flag int) {
	c := colly.NewCollector()
	//c.SetProxy("http://127.0.0.1:1087")
	c.OnRequest(func(r *colly.Request) {
	})
	c.OnError(func(r *colly.Response, err error) {
		r.Request.Ctx.Put("error", err.Error())
		fmt.Errorf("Something went wrong: %v", err)
	})

	if flag == 1 {
		c.OnHTML("div[class]", func(e *colly.HTMLElement) {
			class := e.Attr("class")
			if class == "DetailsHeader-version" {
				//log.Printf("version: %s\n", e.Text)
				e.Response.Ctx.Put("version", e.Text)
				ch <- e.Text
			}
		})
	}

	if flag == 2 {
		c.OnHTML("h1[class]", func(e *colly.HTMLElement) {
			class := e.Attr("class")
			if class == "DetailsHeader-title" {
				//log.Printf("name: %s\n", e.Text)
				if e.Text != "" {
					texts := strings.Split(e.Text, " ")
					e.Response.Ctx.Put("libtype", texts[0])
					//log.Printf("name: %s\n", texts[0])

					ch <- texts[0]
				}
			}
		})
	}

	if flag == 3 {
		c.OnHTML("h2[id]", func(e *colly.HTMLElement) {
			id := e.Attr("id")
			if id == "pkg-overview" {
				printStc := func(firstParagraph string) {
					getFirstSentence := regexp.MustCompile(`\.\s`)
					sentences := getFirstSentence.Split(firstParagraph, -1)
					sentence := strings.Replace(sentences[0], "\n", "", -1)
					//log.Printf("desc: %s\n", sentence)
					e.Response.Ctx.Put("desc", sentence)

					ch <- sentence
				}

				printStc(e.DOM.Next().Text())
			}
		})
	}

	if flag == 4 {
		c.OnHTML("h1[class]", func(e *colly.HTMLElement) {
			class := e.Attr("class")
			if class == "DetailsHeader-title" {
				//log.Printf("name: %s\n", e.Text)
				if e.Text != "" {
					texts := strings.Split(e.Text, " ")
					e.Response.Ctx.Put("pkg name", texts[1])
					//log.Printf("name: %s\n", texts[0])

					ch <- texts[1]
				}
			}
		})
	}

	//c.OnScraped(func(r *colly.Response) {
	//	version := r.Ctx.Get("version")
	//	libtype := r.Ctx.Get("libtype")
	//	desc := r.Ctx.Get("desc")
	//	err := r.Ctx.Get("error")
	//
	//	if err != "" {
	//		chPkg <- &data.Pkg{
	//			Version: "-e-r-r-",
	//		}
	//		return
	//	}
	//
	//	chPkg <- &data.Pkg{
	//		Version:    version,
	//		Libtype:    libtype,
	//		SimpleSpec: desc,
	//	}
	//})

	c.Visit(godevURL + path + "?tab=doc")
}
