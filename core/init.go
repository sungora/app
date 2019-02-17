package core

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// init создание компонента ядра
func init() {
	component = new(componentTyp)
}

var (
	Config    *ConfigRoot   // Корневая конфигурация главного конфигурационного файла
	component *componentTyp // Компонент
)

// компонент
type componentTyp struct {
}

// Init инициализация компонента в приложении
func (comp *componentTyp) Init(cfg *ConfigRoot) (err error) {

	if cfg == nil {
		Config = &ConfigRoot{
			SessionTimeout: 14400,
			TimeZone:       "Europe/Moscow",
			Mode:           "dev",
		}
	} else {
		Config = cfg
	}
	sep := string(os.PathSeparator)

	// техническое имя приложения
	if ext := filepath.Ext(os.Args[0]); ext != "" {
		sl := strings.Split(filepath.Base(os.Args[0]), ext)
		Config.ServiceName = sl[0]
	} else {
		Config.ServiceName = filepath.Base(os.Args[0])
	}

	// читаем конфигурацию
	Config.DirWork, _ = filepath.Abs(filepath.Dir(filepath.Dir(os.Args[0])))
	Config.DirConfig = Config.DirWork + sep + "config"
	Config.DirLog = Config.DirWork + sep + "log"
	Config.DirWww = Config.DirWork + sep + "www"
	path := Config.DirConfig + sep + Config.ServiceName + ".toml"
	if _, err = toml.DecodeFile(path, Config); err != nil {
		return
	}

	// Инициализация временной зоны
	if loc, err := time.LoadLocation(Config.TimeZone); err == nil {
		Config.TimeLocation = loc
	} else {
		Config.TimeLocation = time.UTC
	}

	Config.SessionTimeout *= time.Second

	return
}

// Start запуск компонента в работу
func (comp *componentTyp) Start() (err error) {
	// session
	if 0 < Config.SessionTimeout {
		sessionGC()
	}
	return
}

// Stop завершение работы компонента
func (comp *componentTyp) Stop() (err error) {
	return
}
