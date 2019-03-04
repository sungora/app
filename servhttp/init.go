package servhttp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-chi/chi"

	"github.com/sungora/app/core"
)

// init регистрация компонента в приложении
func init() {
	component = new(componentTyp)
	core.ComponentReg(component)
}

var (
	config    *configFile        // конфигурация
	component *componentTyp      // компонент
	route     = chi.NewRouter(); // роутинг
)

// компонент
type componentTyp struct {
	server    *http.Server  // сервер HTTP
	chControl chan struct{} // управление ожиданием завершения работы сервера
}

// Init инициализация компонента в приложении
func (comp *componentTyp) Init(cfg *core.ConfigRoot) (err error) {

	sep := string(os.PathSeparator)
	config = new(configFile)

	// читаем конфигурацию
	path := cfg.DirConfig + sep + cfg.ServiceName + ".toml"
	if _, err = toml.DecodeFile(path, config); err != nil {
		return
	}

	comp.server = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", config.Http.Host, config.Http.Port),
		Handler:        route,
		ReadTimeout:    time.Second * time.Duration(config.Http.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(config.Http.WriteTimeout),
		IdleTimeout:    time.Second * time.Duration(config.Http.IdleTimeout),
		MaxHeaderBytes: config.Http.MaxHeaderBytes,
	}
	comp.chControl = make(chan struct{})

	return
}

// Start запуск компонента в работу
// Старт сервера HTTP(S)
func (comp *componentTyp) Start() (err error) {
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
func (comp *componentTyp) Stop() (err error) {
	if comp.server == nil {
		return
	}
	if err = comp.server.Shutdown(context.Background()); err != nil {
		return
	}
	<-comp.chControl
	return
}
