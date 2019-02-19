package server

import (
	"context"
	"fmt"
	"github.com/sungora/app/core"
	"github.com/sungora/app/lg"
	"net/http"
	"time"
)

// HandlerFuncHttp Функция обработчик http запросов
type HandlerFuncHttp func(context.Context, *core.RW) (context.Context, *core.RW)

// Служебная структура для построения маршрута и их обработчиков запроса
type routePath struct {
	hp   *HandlerHttp
	path string
}

// Path Формирование машрута
func (rp *routePath) Path(path string) *routePath {
	rp.path += path
	rp.hp.route[rp.path] = make(map[string]HandlerFuncHttp)
	return rp
}

// Use Обработчик конкретного запроса (любой метод) и middleware
func (rp *routePath) Use(hf ...HandlerFuncHttp) *routePath {
	rp.hp.routeMW[rp.path] = append(rp.hp.routeMW[rp.path], hf...)
	return rp
}

// GET Обработчик конкретного запроса и метода GET
func (rp *routePath) GET(hf HandlerFuncHttp) *routePath {
	rp.hp.route[rp.path][http.MethodGet] = hf
	return rp
}

// POST Обработчик конкретного запроса и метода POST
func (rp *routePath) POST(hf HandlerFuncHttp) *routePath {
	rp.hp.route[rp.path][http.MethodPost] = hf
	return rp
}

// PUT Обработчик конкретного запроса и метода PUT
func (rp *routePath) PUT(hf HandlerFuncHttp) *routePath {
	rp.hp.route[rp.path][http.MethodPost] = hf
	return rp
}

// DELETE Обработчик конкретного запроса и метода DELETE
func (rp *routePath) DELETE(hf HandlerFuncHttp) *routePath {
	rp.hp.route[rp.path][http.MethodPost] = hf
	return rp
}

// OPTIONS Обработчик конкретного запроса и метода OPTIONS
func (rp *routePath) OPTIONS(hf HandlerFuncHttp) *routePath {
	rp.hp.route[rp.path][http.MethodPost] = hf
	return rp
}

// //

// Обработчик запросов по протоколу HTTP(S)
type HandlerHttp struct {
	routeMW   map[string][]HandlerFuncHttp          // общие для всех методов и middleware обработчики
	route     map[string]map[string]HandlerFuncHttp // обработчики для конкретного запроса и метода
	cfg       ConfigTyp                             // конфигурация сервера
	server    *http.Server                          // сервер HTTP
	chControl chan struct{}                         // управление ожиданием завершения работы сервера
}

// NewHandlerHttp Конструктор обработчика запросов
func NewHandlerHttp(config ConfigTyp) *HandlerHttp {
	hp := &HandlerHttp{
		routeMW:   make(map[string][]HandlerFuncHttp),
		route:     make(map[string]map[string]HandlerFuncHttp),
		cfg:       config,
		chControl: make(chan struct{}),
	}
	return hp
}

// Path Формирование маршута
func (hp *HandlerHttp) Path(path string) *routePath {
	hp.route[path] = make(map[string]HandlerFuncHttp)
	rp := &routePath{
		hp:   hp,
		path: path,
	}
	return rp
}

// Start Старт сервера HTTP(S)
func (hp *HandlerHttp) Start() (err error) {
	go func() {
		hp.server = &http.Server{
			Addr:           fmt.Sprintf("%s:%d", hp.cfg.Host, hp.cfg.Port),
			Handler:        hp,
			ReadTimeout:    time.Second * time.Duration(hp.cfg.ReadTimeout),
			WriteTimeout:   time.Second * time.Duration(hp.cfg.WriteTimeout),
			IdleTimeout:    time.Second * time.Duration(hp.cfg.IdleTimeout),
			MaxHeaderBytes: hp.cfg.MaxHeaderBytes,
		}
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

func (hp HandlerHttp) String() {
	fmt.Println("middleware")
	for i, v := range hp.routeMW {
		fmt.Printf("%s, %#v\n", i, v)
	}
	fmt.Println("handler")
	for i, v := range hp.route {
		fmt.Printf("%s, %#v\n", i, v)
	}
}

// getHandler Поиск и получение обработчика конкретного запроса и метода
func (hp *HandlerHttp) getHandler(path string, met string) (hf HandlerFuncHttp) {
	if _, ok := hp.route[path][met]; ok {
		return hp.route[path][met]
	}
	return nil
}

// getHandlerMW Поиск и получение общего обработчика конкретного запроса и middleware
func (hp *HandlerHttp) getHandlerMW(path string) (hf []HandlerFuncHttp) {
	if _, ok := hp.routeMW[path]; ok {
		return hp.routeMW[path]
	}
	return nil
}

// ServeHTTP Точка входа запроса (в приложение).
func (hp *HandlerHttp) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var (
		err         error
		handler     = hp.getHandler(r.URL.Path, r.Method)
		handlerList = hp.getHandlerMW(r.URL.Path)
		rw          = core.NewRW(w, r)
	)
	defer r.Body.Close()

	// Обработчики не найдены. Статика.
	if handler == nil && handlerList == nil {
		if err = rw.ResponseStatic(core.Config.DirWww + r.URL.Path); err != nil {
			lg.Error(err)
		}
		return
	}

	// Контекст
	ctx, cancel := context.WithTimeout(context.Background(), hp.server.WriteTimeout)
	defer cancel()

	// middleware
	for i, _ := range handlerList {
		ctx, rw = handlerList[i](ctx, rw)
	}
	if handler != nil {
		ctx, rw = handler(ctx, rw)
	}
}
