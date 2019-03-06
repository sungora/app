package servhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// init регистрация компонента в приложении
// func init() {
// 	component = new(Component)
// 	core.ComponentReg(component)
// }

var (
	config    *Config            // конфигурация
	component *Component         // компонент
	route     = chi.NewRouter(); // роутинг
)

// компонент
type Component struct {
	server    *http.Server  // сервер HTTP
	chControl chan struct{} // управление ожиданием завершения работы сервера
}

// Init инициализация компонента в приложении
func Init(cfg *Config) (com *Component, err error) {
	config = cfg
	component = new(Component)
	component.server = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", config.Http.Host, config.Http.Port),
		Handler:        route,
		ReadTimeout:    time.Second * time.Duration(config.Http.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(config.Http.WriteTimeout),
		IdleTimeout:    time.Second * time.Duration(config.Http.IdleTimeout),
		MaxHeaderBytes: config.Http.MaxHeaderBytes,
	}
	return component, nil
}

// Start запуск компонента в работу
// Старт сервера HTTP(S)
func (comp *Component) Start() (err error) {
	comp.chControl = make(chan struct{})
	go func() {
		if err = comp.server.ListenAndServe(); err != http.ErrServerClosed {
			return
		}
		close(comp.chControl)
	}()
	return
}

// Stop завершение работы компонента
// Остановка сервера HTTP(S)
func (comp *Component) Stop() (err error) {
	if comp.server == nil {
		return
	}
	if err = comp.server.Shutdown(context.Background()); err != nil {
		return
	}
	<-comp.chControl
	return
}
