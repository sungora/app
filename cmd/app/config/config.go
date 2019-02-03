package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"

	"gopkg.in/hostkeybv/app.v1/tool"
)

func init() {
	var c *config
	path := tool.DirConfig + string(os.PathSeparator) + "project.toml"
	if _, err := toml.DecodeFile(path, &c); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
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

var ServiceAccess *serviceaccess
