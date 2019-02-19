package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	errNotFoundHandler = errors.New("handler not found from uri")
)

// HandlerFuncHttp Функция обработчик http запросов
type HandlerFuncHttp func(context.Context)

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
	server    *http.Server                          // сервер HTTP
	chControl chan struct{}                         // управление ожиданием завершения работы сервера
}

// NewHandlerHttp Конструктор обработчика запросов
func NewHandlerHttp(config ConfigTyp) *HandlerHttp {
	hp := &HandlerHttp{
		routeMW:   make(map[string][]HandlerFuncHttp),
		route:     make(map[string]map[string]HandlerFuncHttp),
		chControl: make(chan struct{}),
	}
	hp.server = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler:        hp,
		ReadTimeout:    time.Second * time.Duration(config.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(config.WriteTimeout),
		IdleTimeout:    time.Second * time.Duration(config.IdleTimeout),
		MaxHeaderBytes: config.MaxHeaderBytes,
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

// ServeHTTP Точка входа запроса (в приложение).
func (hp *HandlerHttp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.URL.Path))
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

// getHandler Поиск и получение обработчика конкретного запроса и метода
func (hp *HandlerHttp) getHandler(path string, met string) (hf HandlerFuncHttp, err error) {
	if _, ok := hp.route[path][met]; ok {
		return hp.route[path][met], nil
	}
	return nil, errNotFoundHandler
}

// getHandlerMW Поиск и получение общего обработчика конкретного запроса и middleware
func (hp *HandlerHttp) getHandlerMW(path string) (hf []HandlerFuncHttp, err error) {
	if _, ok := hp.routeMW[path]; ok {
		return hp.routeMW[path], nil
	}
	return nil, errNotFoundHandler
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
