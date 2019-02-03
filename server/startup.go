package server

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/sungora/app/core"
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

	// читаем конфигурацию
	dirWork, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := dirWork + sep + "config" + sep + core.ServiceName + ".toml"
	if _, err = toml.DecodeFile(path, &config); err != nil {
		return
	}

	return
}

// Start запуск компонента в работу
func (comp *componentTyp) Start() (err error) {
	switch config.Server.Proto {
	case "http":
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
	default:
		return errors.New("protocol not defined")
	}
	return
}

// Stop завершение работы компонента
func (comp *componentTyp) Stop() (err error) {
	switch config.Server.Proto {
	case "http":
		if comp.store != nil {
			err = comp.store.Close()
		}
	default:
		return errors.New("protocol not defined")
	}
	return
}
