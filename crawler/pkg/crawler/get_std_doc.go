package crawler

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/huhu/golang/crawler/pkg/data"
)

var defaultURL = "http://localhost:6060/pkg/"

func GetStdLibDoc(baseURL string, lablesMap map[string]data.Datum, pkgsMap map[string]data.Pkg) {
	if baseURL == "" {
		baseURL = defaultURL
	}

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnHTML("td[class]", func(e *colly.HTMLElement) {
		forefather := e.DOM.Parent().Parent().Parent().Parent().Parent().Parent()

		if !forefather.Is("div[id]") {
			log.Fatal("forefather is not div[id],error!")
		}
		stdlib, exist := forefather.Attr("id")
		if !exist {
			log.Println("the id is not exist")
			return
		}
		if stdlib != "stdlib" {
			return
		}
		link, exist := e.DOM.Children().Attr("href")
		if exist {
			fmt.Println(link)

			getDoc(baseURL, link, lablesMap, pkgsMap)
		}
	})

	c.Visit(baseURL)
}

func getDoc(baseURL, link string, lablesMap map[string]data.Datum, pkgsMap map[string]data.Pkg) {
	url := baseURL + link

	c := colly.NewCollector()

	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	path := strings.TrimRight(link, "/")
	//pathes := strings.Split(path, "/")

	p := data.Pkg{
		//Name:     pathes[len(pathes)-1],
		FullPath: path,
	}

	c.OnHTML("h2[class]", func(e *colly.HTMLElement) {
		title := e.Attr("title")
		if title == "Click to hide Overview section" {
			if !e.DOM.Next().Is("p") {
				return
			}
			firstParagraph := e.DOM.Next().Text()
			getFirstSentence := regexp.MustCompile(`\.\s`)
			sentences := getFirstSentence.Split(firstParagraph, -1)
			sentence := strings.Replace(sentences[0], "\n", "", -1)
			p.SimpleSpec = sentence
			pkgsMap[path] = p
			fmt.Printf("\t%s\t%s\n", p.SimpleSpec, p.FullPath)
		}
	})

	//reFunc := regexp.MustCompile(`^(func)`)a
	//reType := regexp.MustCompile(`^(type)`)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)

		//rm3part := regexp.MustCompile(`(^|.)(github.com|golang.org|go.etcd.io|k8s.io|gopkg.in|git.cloud2go.cn).`)
		//if rm3part.MatchString(r.URL.String()) {
		//	fmt.Println("abort!")
		//	r.Abort()
		//}
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnHTML("h3[id], h2[id]", func(e *colly.HTMLElement) {
		id := e.Attr("id")
		datatype := strings.Split(e.DOM.Text(), " ")
		if len(datatype) <= 0 {
			return
		}

		if id == "pkg-variables" {
		}
		switch datatype[0] {
		case "func":
			if !e.DOM.Next().Next().Is("p") {
				return
			}

			var key string
			splitKeys := strings.Split(id, ".")

			if len(splitKeys) > 1 {
				key = splitKeys[len(splitKeys)-1]
			} else {
				key = splitKeys[0]
			}
			firstParagraph := e.DOM.Next().Next().Text()
			getFirstSentence := regexp.MustCompile(`\.\s`)
			sentences := getFirstSentence.Split(firstParagraph, -1)
			sentence := strings.Replace(sentences[0], "\n", "", -1)

			i := data.Datum{
				Name:       key,
				Definition: e.DOM.Next().Text(),
				Label:      id,
				Pkg:        strings.TrimRight(link, "/"),
				Datatype:   "func",
				SimpleSpec: sentence,
			}
			lablesMap[link+id] = i
			//fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\n", i.name, i.datatype, i.label, i.definition, i.simpleSpec, i.pkg)

		case "type":
			//fmt.Println("type", e.DOM.Next().Is("p"))

			var sentence, dataType, definition string
			next := e.DOM.Next()
			for i := 0; !next.Is("h2, h3, div"); i++ {
				if i == 0 && next.Is("p") {
					// first paragraph
					firstParagraph := next.Text()
					getFirstSentence := regexp.MustCompile(`\.\s`)
					sentences := getFirstSentence.Split(firstParagraph, -1)
					sentence = strings.Replace(sentences[0], "\n", "", -1)
				}
				if next.Is("pre") {
					// is code
					definitions := strings.Split(next.Text(), " ")
					if len(definitions) <= 0 {
						continue
					}
					if definitions[0] == "type" {
						// is definition
						switch definitions[2] {
						case "struct":
							dataType = "struct"
						case "interface":
							dataType = "interface"
						default:
							dataType = "other"
						}
						definition = next.Text()
					}
				}

				next = next.Next()
			}

			var key string
			splitKeys := strings.Split(id, ".")

			if len(splitKeys) > 1 {
				key = splitKeys[len(splitKeys)-1]
			} else {
				key = splitKeys[0]
			}

			i := data.Datum{
				Name:       key,
				Definition: definition,
				Label:      id,
				Pkg:        strings.TrimRight(link, "/"),
				Datatype:   dataType,
				SimpleSpec: sentence,
			}

			lablesMap[link+id] = i

			//fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\n", i.name, i.datatype, i.label, i.definition, i.simpleSpec, i.pkg)
			//time.Sleep(1 * time.Second)
		}

	})

	c.Visit(url)
}
