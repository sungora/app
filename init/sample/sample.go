package sample

import (
	"flag"
	"fmt"
	"os"
	"time"

	middlewareChi "github.com/go-chi/chi/middleware"

	"github.com/sungora/app"
	"github.com/sungora/app/connect"
	"github.com/sungora/app/core"
	"github.com/sungora/app/lg"
	"github.com/sungora/app/servhttp"
	"github.com/sungora/app/servhttp/middlew"
	"github.com/sungora/app/workflow"

	"github.com/sungora/app/internal/sample"
	"github.com/sungora/app/internal/sample/config"
)

type Config struct {
	Core     core.Config     `yaml:"Core"`
	Lg       lg.Config       `yaml:"Lg"`
	Workflow workflow.Config `yaml:"Workflow"`
	Http     servhttp.Config `yaml:"Http"`
	Connect  connect.Config  `yaml:"Connect"`
	Sample   config.Config   `yaml:"Calculator"`
}

const (
	version = "1.0.0"
)

func Init() (code int) {
	var (
		err             error
		component       app.Componenter
		componentServer *servhttp.Component
	)

	flagConfigPath := flag.String("c", "config/sample.yaml", "used for set path to config file")
	flag.Parse()


	lg.Dumper(os.Getwd())

	// загрузка конфигурации
	cfg := &Config{}
	if err = app.LoadConfig(*flagConfigPath, cfg); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1
	}

	// COMPONENTS
	// core
	if component, err = core.Init(&cfg.Core, version); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1
	}
	app.ComponentAdd(component)
	// logs
	if component, err = lg.Init(&cfg.Lg, cfg.Core.ServiceName); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1
	}
	app.ComponentAdd(component)
	// connect
	if component, err = connect.Init(&cfg.Connect); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1
	}
	app.ComponentAdd(component)
	// workflow
	if component, err = workflow.Init(&cfg.Workflow); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1
	}
	app.ComponentAdd(component)
	// servhttp
	if componentServer, err = servhttp.Init(&cfg.Http); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1
	}
	app.ComponentAdd(componentServer)
	route := componentServer.GetRoute()

	// APPLICATION
	route.NotFound(middlew.NotFound)
	route.Use(middlew.TimeoutContext(time.Second * time.Duration(cfg.Http.WriteTimeout-1)))
	route.Use(middlewareChi.Recoverer)
	route.Use(middlewareChi.Logger)

	// MODULES
	if err = sample.Init(&cfg.Sample); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1
	}

	// запуск и остановка приложения
	var isStart = int8(2)
	return app.Start(&isStart)
}
