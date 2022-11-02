package controller

import "testing"

func TestCtlCacheMeta_Produce(t *testing.T) {
	//dbPath := "../svc/test.db"

	svc := GetCtlCacheMeta()

	svc.Produce()
}
