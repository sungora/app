package core

import "time"

// главная конфигурация
type configMain struct {
	SessionTimeout time.Duration
	TimeZone       string
	UseDB          bool
	Mode           string
}
