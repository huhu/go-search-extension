package svc

import "testing"

func TestGetAndSaveMetasFromFile(t *testing.T) {
	t.Run("GetAndSaveMetasFromFile", func(t *testing.T) {
		GetAndSaveMetasFromFile("./test.db", "/Users/cym/tmp/packages")
	})
}
