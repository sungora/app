package tool

import (
	"time"
)

const (
	DATETIME = "2006-01-02 15:04:05"
)

var (
	DirConfig    string         //
	DirWork      string         //
	DirLog       string         //
	DirWww       string         //
	DirTpl       string         //
	DirPid       string         //
	ServiceName  string         // Техническое название приложения
	TimeLocation *time.Location // Временная зона
)
