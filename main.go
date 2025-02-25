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

	//spark.SendSparkHttp(context.Background(), []float64{17.06, 16.95, 16.81, 16.78, 16.95, 17.02,
	//	17.08, 17.06, 17.26, 16.75, 16.35, 17.17, 17.28, 17.22, 17.31, 16.82, 18.15,
	//	18.13, 18.15, 18.43, 18.05, 18.01, 17.70, 18.24, 17.98, 19.02, 18.88, 18.39, 17.02, 16.48}, "", "1")
}
