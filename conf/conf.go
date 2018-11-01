package conf

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

type ConfigMain struct {
	Name           string
	TimeZone       string
	DriverDB       string // Драйвер DB
	SessionTimeout time.Duration
	Host           string
	Port           int
	Mode           string // Режим работы приложения
}

var DirWork string
var DirConfig string
var DirStatic string
var TimeLocation *time.Location
