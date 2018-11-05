package main

import (
	"os"

	_ "PKGAPPNAME/config"        // конфигурация
	_ "PKGAPPNAME/config/lg"     // log
	_ "PKGAPPNAME/config/route"  // route controllers
	_ "PKGAPPNAME/config/worker" // workers task cron

	"gopkg.in/sungora/app.v1/server"
)

func main() {

	os.Exit(server.Start())

}
