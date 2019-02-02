package setup

import (
	"accounter/controller/service"
	"accounter/controller/start"

	"github.com/hostkeybv/app"
)

func setupController() {

	app.Route.Set("/", start.NewControlMain)
	app.Route.Set("/service/control", service.NewControlService)

}
