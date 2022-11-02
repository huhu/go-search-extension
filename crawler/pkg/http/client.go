package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/huhu/golang/crawler/pkg/data"
)

var (
	docURL   = "api.godoc.org"
	queryURL = "https://api.godoc.org/search?q=\"%s\""
	cacheURL = "https://index.golang.org/index?since=%s&limit=%d"
)

// CliHTTP http client
type CliHTTP struct {
	*http.Client
}

// NewCliHTTP new a http client
func NewCliHTTP(useProxy bool) *CliHTTP {
	var transport http.RoundTripper
	if useProxy {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse("http://127.0.0.1:1087")
		}

		transport = &http.Transport{Proxy: proxy}
	}
	return &CliHTTP{
		&http.Client{
			Transport: transport,
		},
	}
}

// GetMetadatas get metadatas of the pachage by path
func (cli *CliHTTP) GetMetadatas(path string) []data.Metadata {
	query := fmt.Sprintf(queryURL, path)

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		fmt.Errorf("new req failed due to %s", err)
		return nil
	}
	req.Header.Set("Authority", "api.godoc.org")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

	resp, err := cli.Do(req)
	if err != nil {
		fmt.Errorf("http req error in getSpec of %s: %v", path, err)
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read resp body in getSpec: %v \n", err)
		return nil
	}
	meta := &data.Result{}

	json.Unmarshal(
		body, meta,
	)
	return meta.Meta
}

type Meta struct {
	Path      string `json:"Path"`
	Version   string `json:"Version"`
	Timestamp string `json:"Timestamp"`
}

func (cli *CliHTTP) GetMetadatasFromCache(since string, limit int) ([]Meta, error) {
	query := fmt.Sprintf(cacheURL, since, limit)

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		fmt.Errorf("new req failed due to %s", err)
		return []Meta{}, err
	}
	resp, err := cli.Do(req)
	if err != nil {
		fmt.Errorf("http req error of since %s and limit %d: %v", since, limit, err)
		return []Meta{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read resp body in getSpec: %v \n", err)
		return []Meta{}, err
	}

	//log.Printf("body : %v\n", body)

	metas := make([]Meta, 0, limit)
	var startIndex int
	for i, char := range body {
		if char == 10 {
			meta := &Meta{}
			err := json.Unmarshal(
				body[startIndex:i], meta,
			)

			//log.Printf("i: %d; meta json: %v\n", i, meta)

			if err != nil {
				return []Meta{}, err
			}
			metas = append(metas, *meta)
			startIndex = i + 1
		}
	}

	//log.Printf("body json: %v\n", metas)

	return metas, nil
}
