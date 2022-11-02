package main

import "github.com/huhu/golang/crawler/pkg/controller"

func main() {
	dbPath := "./pkg/controller/test.db"

	ctl := controller.GetCtlUpdate(dbPath)
	ctl.Start()
}
