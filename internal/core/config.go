package core

import (
	"github.com/sungora/app"
	"github.com/sungora/app/connect"
	"github.com/sungora/app/lg"
	"github.com/sungora/app/servhttp"
	"github.com/sungora/app/workflow"
)

type Config struct {
	App      app.Config      `yaml:"App"`
	Lg       lg.Config       `yaml:"Lg"`
	Workflow workflow.Config `yaml:"Workflow"`
	Http     servhttp.Config `yaml:"Http"`
	Connect  connect.Config  `yaml:"Connect"`
	Sample   sample          `yaml:"Sample"`
}

type sample struct {
	Name     string  `yaml:"Name"`
	IsAccess bool    `yaml:"IsAccess"`
	Balance  float64 `yaml:"Balance"`
}

var Cfg = &Config{}
