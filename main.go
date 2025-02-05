package main

import (
	_ "financia/public/db/connector"
	"financia/router"
	"financia/server"
)

func main() {
	go server.CronDailyWorker()
	router.HTTPRouter()
}
