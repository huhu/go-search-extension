package http

import (
	"log"
	"testing"
)

func TestCliHTTP_GetMetadatasFromCache(t *testing.T) {
	//since := "2019-04-10T19:08:52.997264Z"
	limit := 10

	since2 := "2020-06-10T19:08:52.997264Z"

	t.Run("GetMetadataFromCache", func(t *testing.T) {
		cli := NewCliHTTP(true)
		metas, err := cli.GetMetadatasFromCache(since2, limit)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("body json: %v\n", metas)

	})
}
