package core

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	Cfg       *Config    // Корневая конфигурация главного конфигурационного файла
	component *Component // Компонент
)

// компонент
type Component struct {
}

// Init инициализация компонента в приложении
func Init(cfg *Config) (com *Component, err error) {

	Cfg = cfg
	component = new(Component)

	// временная зона
	if Cfg.TimeZone != "" {
		Cfg.TimeZone = "Europe/Moscow"
	}
	if loc, err := time.LoadLocation(Cfg.TimeZone); err == nil {
		Cfg.TimeLocation = loc
	} else {
		Cfg.TimeLocation = time.UTC
	}
	// режим работы приложения
	if Cfg.TimeZone != "" {
		Cfg.Mode = "dev"
	}
	// техническое имя приложения
	if Cfg.ServiceName != "" {
		if ext := filepath.Ext(os.Args[0]); ext != "" {
			sl := strings.Split(filepath.Base(os.Args[0]), ext)
			Cfg.ServiceName = sl[0]
		} else {
			Cfg.ServiceName = filepath.Base(os.Args[0])
		}
	}
	// пути
	sep := string(os.PathSeparator)
	if Cfg.DirWork == "" {
		Cfg.DirWork = filepath.Dir(filepath.Dir(os.Args[0]))
		if Cfg.DirWork == "." {
			Cfg.DirWork, _ = os.Getwd()
			Cfg.DirWork = filepath.Dir(Cfg.DirWork)
		}
	}
	if Cfg.DirConfig == "" {
		Cfg.DirConfig = Cfg.DirWork + sep + "config"
	}
	if Cfg.DirLog == "" {
		Cfg.DirLog = Cfg.DirWork + sep + "log"
	}
	if Cfg.DirWww == "" {
		Cfg.DirWww = Cfg.DirWork + sep + "www"
	}
	// сессия
	Cfg.SessionTimeout *= time.Second

	return component, nil
}

// Start запуск компонента в работу
func (comp *Component) Start() (err error) {
	// session
	if 0 < Cfg.SessionTimeout {
		sessionGC()
	}
	return
}

// Stop завершение работы компонента
func (comp *Component) Stop() (err error) {
	return
}
