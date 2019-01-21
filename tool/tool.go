package tool

import (
	"time"
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
