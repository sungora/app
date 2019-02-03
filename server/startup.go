package server

import (
	"fmt"
	"net"
	"net/http"
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
	store net.Listener
}

var (
	config    *configMain   // конфигурация
	component *componentTyp // компонент
)

// Init инициализация компонента в приложении
func (comp *componentTyp) Init() (err error) {
	sep := string(os.PathSeparator)
	config = new(configMain)

	// техническое имя приложения
	if ext := filepath.Ext(os.Args[0]); ext != "" {
		sl := strings.Split(filepath.Base(os.Args[0]), ext)
		config.ServiceName = sl[0]
	} else {
		config.ServiceName = filepath.Base(os.Args[0])
	}

	// читаем конфигурацию
	dirWork, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := dirWork + sep + "config" + sep + config.ServiceName + ".toml"
	if _, err = toml.DecodeFile(path, &config); err != nil {
		return
	}

	return
}

// Start запуск компонента в работу
func (comp *componentTyp) Start() (err error) {
	Server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port),
		Handler:        new(serverHttp),
		ReadTimeout:    time.Second * time.Duration(config.Server.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(config.Server.WriteTimeout),
		MaxHeaderBytes: config.Server.MaxHeaderBytes,
	}
	if comp.store, err = net.Listen("tcp", Server.Addr); err != nil {
		return
	}
	go Server.Serve(comp.store)
	return
}

// Stop завершение работы компонента
func (comp *componentTyp) Stop() (err error) {
	err = comp.store.Close()
	return
}
