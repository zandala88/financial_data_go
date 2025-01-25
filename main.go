package main

import (
	_ "financia/public/db/connector"
	"financia/router"
)

func main() {
	router.HTTPRouter()
}
