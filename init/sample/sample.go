package sample

import (
	"flag"
	"fmt"
	"os"

	"github.com/sungora/app"
	"github.com/sungora/app/connect"
	"github.com/sungora/app/internal/core"
	"github.com/sungora/app/lg"
	"github.com/sungora/app/servhttp"
	"github.com/sungora/app/workflow"
)

func Start() (code int) {
	var (
		err             error
		component       app.Componenter
		componentServer *servhttp.Component
	)

	// Флаги
	flagConfigPath := flag.String("c", "config/sample.yaml", "used for set path to config file")
	flag.Parse()

	// загрузка конфигурации
	if err = app.LoadConfig(*flagConfigPath, core.Cfg); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1
	}

	// COMPONENTS
	// logs
	if component, err = lg.Init(&core.Cfg.Lg); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1
	}
	app.ComponentAdd(component)
	// connect
	if component, err = connect.Init(&core.Cfg.Connect); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1
	}
	app.ComponentAdd(component)
	// workflow
	if component, err = workflow.Init(&core.Cfg.Workflow); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1
	}
	app.ComponentAdd(component)
	// servhttp
	if componentServer, err = servhttp.Init(&core.Cfg.Http); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1
	}
	app.ComponentAdd(componentServer)

	// APPLICATION
	// routes
	routes(componentServer.GetRoute())
	// workers
	workers()
	// logs
	logs()

	// START запуск и остановка приложения
	if err = app.StartLock(&core.Cfg.App); err != nil {
		return 1
	}
	return
}
