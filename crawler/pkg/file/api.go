package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/huhu/golang/crawler/pkg/data"
)

func GetDataFromFile(filePath string) *data.Result {
	JSONParse := &data.Result{}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("read file err: %s", err)
	}

	err = json.Unmarshal(data, JSONParse)
	if err != nil {
		log.Fatalf("unmarshal err: %s", err)
	}

	fmt.Printf("len of metadatas: %d\n", len(JSONParse.Meta))
	return JSONParse
}
