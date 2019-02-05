package core

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// init создание компонента ядра
func init() {
	component = new(componentTyp)
}

// компонент
type componentTyp struct {
}

var (
	Config    *ConfigRoot   // Корневая конфигурация главного конфигурационного файла
	component *componentTyp // Компонент
)

// Init инициализация компонента в приложении
func (comp *componentTyp) Init(cfg *ConfigRoot) (err error) {

	if cfg == nil {
		Config = new(ConfigRoot)
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

// ///////////////

// Интерфейс компонентов приложения
type Componenter interface {
	// Инициализация компонентов приложения
	Init(cfg *ConfigRoot) (err error)
	// Запуск в работу компонентов приложения
	Start() (err error)
	// Завершение работы компонентов приложения
	Stop() (err error)
}

// Срез зарегитрированных компонентов приложения
var componentList []Componenter

// ComponentReg Регистрация компонента приложения
func ComponentReg(com Componenter) {
	componentList = append(componentList, com)
}

var (
	// Канал управления запуском и остановкой приложения
	chanelAppControl = make(chan os.Signal, 1)
)

// Start Launch an application
func Start() (code int) {

	defer func() {
		chanelAppControl <- os.Interrupt
	}()
	var err error
	componentList = append([]Componenter{component}, componentList...)

	if len(componentList) == 1 {
		fmt.Fprintln(os.Stderr, "Ни одного компонента не зарегистрировано")
		return 1
	}

	// инициализация компонентов
	for i := 0; i < len(componentList); i++ {
		fmt.Fprintf(os.Stdout, "Init component %s\n", packageName(componentList[i]))
		if err = componentList[i].Init(Config); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
	}

	// 	завершение работы компонентов
	defer func() {
		for i := 0; i < len(componentList); i++ {
			fmt.Fprintf(os.Stdout, "Stop component %s\n", packageName(componentList[i]))
			if err = componentList[i].Stop(); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				code = 1
			}
		}
	}()

	// начало в работы компонентов
	for i := 0; i < len(componentList); i++ {
		fmt.Fprintf(os.Stdout, "Start component %s\n", packageName(componentList[i]))
		if err = componentList[i].Start(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
	}

	// The correctness of the application is closed by a signal
	signal.Notify(chanelAppControl, os.Interrupt)
	<-chanelAppControl
	return
}

// Stop an application
func Stop() {
	chanelAppControl <- os.Interrupt
	<-chanelAppControl
}

// packageName Получение уникального имени пакета
func packageName(obj interface{}) string {
	var rt = reflect.TypeOf(obj)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return rt.PkgPath()
}
