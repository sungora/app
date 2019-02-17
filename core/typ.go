package core

import "time"

// главная конфигурация
type ConfigRoot struct {
	SessionTimeout time.Duration  //
	TimeZone       string         //
	Mode           string         //
	DirWork        string         //
	DirConfig      string         //
	DirLog         string         //
	DirWww         string         //
	ServiceName    string         // Техническое название приложения
	TimeLocation   *time.Location // Временная зона
}
