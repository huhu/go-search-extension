package svc

import "testing"

func TestGetAndSaveStdDoc(t *testing.T) {
	t.Run("GetAndSaveStdDoc", func(t *testing.T) {
		GetAndSaveStdDoc("./test.db", "")
	})
}
