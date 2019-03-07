package core

import (
	"time"
)

// главная конфигурация
type Config struct {
	SessionTimeout time.Duration  `yaml:"SessionTimeout"` //
	TimeZone       string         `yaml:"TimeZone"`       //
	Mode           string         `yaml:"Mode"`           //
	DirWork        string         `yaml:"DirWork"`        //
	DirConfig      string         `yaml:"DirConfig"`      //
	DirLog         string         `yaml:"DirLog"`         //
	DirWww         string         `yaml:"DirWww"`         //
	ServiceName    string         `yaml:"ServiceName"`    // Техническое название приложения
	TimeLocation   *time.Location ``                      // Временная зона
}
