package main

import (
	_ "financia/public/db/connector"
	"financia/router"
	"financia/server"
	"financia/server/python"
)

func main() {
	go python.NewGRPCClient()
	go server.CronDailyWorker()
	router.HTTPRouter()
}
