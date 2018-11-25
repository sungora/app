package tool

import (
	"os"
	"time"
)

func init() {
	DirWork, _ = os.Getwd()
	sep := string(os.PathSeparator)
	DirConfig = DirWork + sep + "config"
	DirStatic = DirWork + sep + "www"
}

var DirWork string
var DirConfig string
var DirStatic string
var TimeLocation *time.Location
