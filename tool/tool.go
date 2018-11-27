package tool

import (
	"os"
	"time"
)

func init() {
	DirWork, _ = os.Getwd()
	Sep := string(os.PathSeparator)
	DirConfig = DirWork + Sep + "config"
	DirWww = DirWork + Sep + "www"
	DirTpl = DirWork + Sep + "tpl"
}

var (
	DirWork      string
	DirConfig    string
	DirWww       string
	DirTpl       string
	TimeLocation *time.Location
)
