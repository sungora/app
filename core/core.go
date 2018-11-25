package core

import (
	"time"
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