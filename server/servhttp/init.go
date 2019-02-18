package servhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/sungora/app/core"
)

// init регистрация компонента в приложении
func init() {
	component = new(componentTyp)
	core.ComponentReg(component)
}

var (
	config          *configMain   // конфигурация
	component       *componentTyp // компонент
	ErrorNotDefinde = errors.New("protocol not defined")
)

const (
	proto_http  = "http"
	proto_https = "https"
)

// компонент
type componentTyp struct {
	ser       *http.Server
	chControl chan struct{} // канал управление ожиданием завершения
}

// Init инициализация компонента в приложении
func (comp *componentTyp) Init(cfg *core.ConfigRoot) (err error) {
	sep := string(os.PathSeparator)
	config = new(configMain)

	// читаем конфигурацию
	path := cfg.DirConfig + sep + cfg.ServiceName + ".toml"
	if _, err = toml.DecodeFile(path, config); err != nil {
		return
	}

	// Инициализируем сервер
	switch config.Server.Proto {
	case proto_http:
		comp.ser = &http.Server{
			Addr:           fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port),
			Handler:        new(serverHttp),
			ReadTimeout:    time.Second * time.Duration(config.Server.ReadTimeout),
			WriteTimeout:   time.Second * time.Duration(config.Server.WriteTimeout),
			IdleTimeout:    time.Second * time.Duration(config.Server.IdleTimeout),
			MaxHeaderBytes: config.Server.MaxHeaderBytes,
		}
	default:
		return ErrorNotDefinde
	}

	// управление ожиданием завершения
	comp.chControl = make(chan struct{})
	return
}

// Start запуск компонента в работу
func (comp *componentTyp) Start() (err error) {
	switch config.Server.Proto {
	case proto_http:
		go func() {
			if err = comp.ser.ListenAndServe(); err != http.ErrServerClosed {
				return
			}
			close(comp.chControl)
		}()
	default:
		return ErrorNotDefinde
	}
	return
}

// Stop завершение работы компонента
func (comp *componentTyp) Stop() (err error) {
	switch config.Server.Proto {
	case proto_http:
		if comp.ser != nil {
			if err = comp.ser.Shutdown(context.Background()); err != nil {
				return
			}
			<-comp.chControl
		}
	default:
		return ErrorNotDefinde
	}
	return
}
