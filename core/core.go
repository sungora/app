package core

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ConfigTyp struct {
	Name           string
	TimeZone       string
	DriverDB       string // Драйвер DB
	SessionTimeout time.Duration
	Host           string
	Port           int
	Mode           string // Режим работы приложения
}

var Config *ConfigTyp
var DB *gorm.DB
