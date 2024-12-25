package main

import (
	"financia/router"
	"financia/server"
)

func main() {
	go server.InsertDailyDate()
	router.HTTPRouter()
}
