package servhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

var (
	config    *Config    // конфигурация
	component *Component // компонент
)

// Init инициализация компонента в приложении
func Init(cfg *Config) (com *Component, router *chi.Mux, err error) {
	config = cfg
	component = &Component{
		Server: &http.Server{
			Addr:           fmt.Sprintf("%s:%d", config.Host, config.Port),
			Handler:        chi.NewRouter(),
			ReadTimeout:    time.Second * config.ReadTimeout,
			WriteTimeout:   time.Second * config.WriteTimeout,
			IdleTimeout:    time.Second * config.IdleTimeout,
			MaxHeaderBytes: config.MaxHeaderBytes,
		},
	}
	return component, component.Server.Handler.(*chi.Mux), nil
}

// Start запуск компонента в работу
// Старт сервера HTTP(S)
func (comp *Component) Start() (err error) {
	comp.chControl = make(chan struct{})
	go func() {
		if err = comp.Server.ListenAndServe(); err != http.ErrServerClosed {
			return
		}
		close(comp.chControl)
	}()
	return
}

// Stop завершение работы компонента
// Остановка сервера HTTP(S)
func (comp *Component) Stop() (err error) {
	if comp.Server == nil {
		return
	}
	if err = comp.Server.Shutdown(context.Background()); err != nil {
		return
	}
	<-comp.chControl
	return
}

// GetRoute получение обработчика запросов
func GetRoute() *chi.Mux {
	return component.Server.Handler.(*chi.Mux)
}

func GetConfig() *Config {
	return config
}
