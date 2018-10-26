package conf

import (
	"os"
	"time"
)

func init() {
	dir, _ := os.Getwd()
	sep := string(os.PathSeparator)
	ConfigDir = dir + sep + "config" + sep
}

var ConfigDir string
var TimeLocation *time.Location
