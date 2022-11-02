package crawler

import (
	"log"
	"strings"
	"testing"
)

func TestReg(t *testing.T) {
	text := "package hello"
	texts := strings.Split(text, " ")
	log.Println(texts[len(texts)-1])

	ch := make(chan string, 1)

	GetDataFromPkgGo("github.com/hawkular/hawkular-client-go", ch, 3)
	r := <-ch
	log.Println(r)
}
