package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/sungora/app.v1/conf"
)

func init() {
	var c *config
	path := conf.DirConfig + string(os.PathSeparator) + "project.toml"
	if _, err := toml.DecodeFile(path, &c); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	ServiceAccess = c.ServiceAccess
}

type config struct {
	ServiceAccess *serviceaccess
}

type serviceaccess struct {
	UrlApi   string
	Login    string
	Password string
}

var ServiceAccess = new(serviceaccess)
