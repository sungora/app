package sample

import (
	"flag"
	"fmt"
	"os"
	"time"

	middlewareChi "github.com/go-chi/chi/middleware"

	"github.com/sungora/app"
	"github.com/sungora/app/connect"
	"github.com/sungora/app/lg"
	"github.com/sungora/app/servhttp"
	"github.com/sungora/app/servhttp/middlew"
	"github.com/sungora/app/workflow"

	"github.com/sungora/app/internal/sample"
	"github.com/sungora/app/internal/sample/config"
)

type Config struct {
	App      app.Config      `yaml:"App"`
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

	// Флаги
	flagConfigPath := flag.String("c", "config/sample.yaml", "used for set path to config file")
	flag.Parse()

	// загрузка конфигурации
	cfg := &Config{}
	if err = app.LoadConfig(*flagConfigPath, cfg); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1
	}

	// COMPONENTS
	// logs
	if component, err = lg.Init(&cfg.Lg); err != nil {
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
	if err = app.StartLock(&cfg.App); err != nil {
		return 1
	}
	return
}
