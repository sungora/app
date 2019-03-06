package core

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// init создание компонента ядра
func init() {
	component = new(Component)
}

var (
	config    *Config    // Корневая конфигурация главного конфигурационного файла
	component *Component // Компонент
)

// компонент
type Component struct {
}

// Init инициализация компонента в приложении
func Init(cfg *Config) (com *Component, err error) {

	config = cfg
	component = new(Component)

	// временная зона
	if config.TimeZone != "" {
		config.TimeZone = "Europe/Moscow"
	}
	if loc, err := time.LoadLocation(config.TimeZone); err == nil {
		config.TimeLocation = loc
	} else {
		config.TimeLocation = time.UTC
	}
	// режим работы приложения
	if config.TimeZone != "" {
		config.Mode = "dev"
	}
	// техническое имя приложения
	if config.ServiceName != "" {
		if ext := filepath.Ext(os.Args[0]); ext != "" {
			sl := strings.Split(filepath.Base(os.Args[0]), ext)
			config.ServiceName = sl[0]
		} else {
			config.ServiceName = filepath.Base(os.Args[0])
		}
	}
	// пути
	sep := string(os.PathSeparator)
	if config.DirWork == "" {
		config.DirWork, _ = filepath.Abs(filepath.Dir(filepath.Dir(os.Args[0])))
	}
	if config.DirConfig == "" {
		config.DirConfig = config.DirWork + sep + "config"
	}
	if config.DirLog == "" {
		config.DirLog = config.DirWork + sep + "log"
	}
	if config.DirWww == "" {
		config.DirWww = config.DirWork + sep + "www"
	}
	// сессия
	config.SessionTimeout *= time.Second

	return component, nil
}

// Start запуск компонента в работу
func (comp *Component) Start() (err error) {
	// session
	if 0 < config.SessionTimeout {
		sessionGC()
	}
	return
}

// Stop завершение работы компонента
func (comp *Component) Stop() (err error) {
	return
}

func GetConfig() {

}
