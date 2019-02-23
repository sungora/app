package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// Обработчик запросов по протоколу HTTP(S)
type HandlerHttp struct {
	server    *http.Server                          // сервер HTTP
	chControl chan struct{}                         // управление ожиданием завершения работы сервера
}

// NewHandlerHttp Конструктор обработчика запросов
func NewHandlerHttp(cfg ConfigTyp, mux *chi.Mux) *HandlerHttp {
	hp := &HandlerHttp{
		chControl: make(chan struct{}),
		server: &http.Server{
			Addr:           fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler:        mux,
			ReadTimeout:    time.Second * time.Duration(cfg.ReadTimeout),
			WriteTimeout:   time.Second * time.Duration(cfg.WriteTimeout),
			IdleTimeout:    time.Second * time.Duration(cfg.IdleTimeout),
			MaxHeaderBytes: cfg.MaxHeaderBytes,
		},
	}
	return hp
}


// Start Старт сервера HTTP(S)
func (hp *HandlerHttp) Start() (err error) {
	go func() {
		if err = hp.server.ListenAndServe(); err != http.ErrServerClosed {
			return
		}
		close(hp.chControl)
	}()
	return
}

// Stop Остановка сервера HTTP(S)
func (hp *HandlerHttp) Stop() (err error) {
	if hp.server == nil {
		return
	}
	if err = hp.server.Shutdown(context.Background()); err != nil {
		return
	}
	<-hp.chControl
	return
}
