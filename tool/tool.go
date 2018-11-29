package tool

import (
	"os"
	"time"
)

func init() {
	DirWork, _ = os.Getwd()
	sep := string(os.PathSeparator)
	DirConfig = DirWork + sep + "config"
	DirWww = DirWork + sep + "www"
	DirTpl = DirWork + sep + "tpl"
}

var (
	DirWork      string
	DirConfig    string
	DirWww       string
	DirTpl       string
	TimeLocation *time.Location
)
