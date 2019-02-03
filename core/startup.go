package core

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/sungora/app/startup"
)

// init регистрация компонента в приложении
func init() {
	component = new(componentTyp)
	startup.SetComponent(component)
}

// компонент
type componentTyp struct {
}

var (
	Config       *configMain    // конфигурация
	component    *componentTyp  // компонент
	DirWork      string         //
	DirConfig    string         //
	DirLog       string         //
	DirWww       string         //
	ServiceName  string         // Техническое название приложения
	TimeLocation *time.Location // Временная зона
)

// Init инициализация компонента в приложении
func (comp *componentTyp) Init() (err error) {

	sep := string(os.PathSeparator)
	Config = new(configMain)

	// техническое имя приложения
	if ext := filepath.Ext(os.Args[0]); ext != "" {
		sl := strings.Split(filepath.Base(os.Args[0]), ext)
		ServiceName = sl[0]
	} else {
		ServiceName = filepath.Base(os.Args[0])
	}

	// читаем конфигурацию
	DirWork, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	DirConfig = DirWork + sep + "config"
	DirLog = DirWork + sep + "log"
	DirWww = DirWork + sep + "www"
	path := DirConfig + sep + ServiceName + ".toml"
	if _, err = toml.DecodeFile(path, &Config); err != nil {
		return
	}

	// Инициализация временной зоны
	if loc, err := time.LoadLocation(Config.TimeZone); err == nil {
		TimeLocation = loc
	} else {
		TimeLocation = time.UTC
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
