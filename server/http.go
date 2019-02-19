package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sungora/app/core"
	"github.com/sungora/app/lg"
)

// HandlerFuncHttp Функция обработчик http запросов
type HandlerFuncHttp func(context.Context, *core.RW) (context.Context, *core.RW)

// Служебная структура для построения маршрута и их обработчиков запроса
type routePath struct {
	hp   *HandlerHttp
	path string
}

// Route Формирование машрута
func (rp *routePath) Route(p string) *routePath {
	rp.hp.route[rp.path+p] = make(map[string]HandlerFuncHttp)
	rpNew := &routePath{
		hp:   rp.hp,
		path: rp.path + p,
	}
	fmt.Println(rp.path)
	return rpNew
}

// Before Общие обработчики запросов (middleware)
func (rp *routePath) Before(hf ...HandlerFuncHttp) *routePath {
	rp.hp.routeBe[rp.path] = append(rp.hp.routeBe[rp.path], hf...)
	return rp
}

// After Общие обработчики запросов (middleware)
func (rp *routePath) After(hf ...HandlerFuncHttp) *routePath {
	rp.hp.routeAf[rp.path] = append(rp.hp.routeAf[rp.path], hf...)
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
	routeBe   map[string][]HandlerFuncHttp          // общие обработчики запросов (middleware)
	routeAf   map[string][]HandlerFuncHttp          // общие обработчики запросов (middleware)
	route     map[string]map[string]HandlerFuncHttp // обработчики для конкретного запроса и метода
	server    *http.Server                          // сервер HTTP
	chControl chan struct{}                         // управление ожиданием завершения работы сервера
}

// NewHandlerHttp Конструктор обработчика запросов
func NewHandlerHttp(cfg ConfigTyp) *HandlerHttp {
	hp := &HandlerHttp{
		routeBe:   make(map[string][]HandlerFuncHttp),
		routeAf:   make(map[string][]HandlerFuncHttp),
		route:     make(map[string]map[string]HandlerFuncHttp),
		chControl: make(chan struct{}),
	}
	hp.server = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:        hp,
		ReadTimeout:    time.Second * time.Duration(cfg.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(cfg.WriteTimeout),
		IdleTimeout:    time.Second * time.Duration(cfg.IdleTimeout),
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}
	return hp
}

// Route Формирование маршута
func (hp *HandlerHttp) Route(p string) *routePath {
	hp.route[p] = make(map[string]HandlerFuncHttp)
	rp := &routePath{
		hp:   hp,
		path: p,
	}
	return rp
}

// getHandler Поиск и получение обработчика конкретного запроса и метода
func (hp *HandlerHttp) getHandler(path string, met string) (hf HandlerFuncHttp) {
	if _, ok := hp.route[path][met]; ok {
		return hp.route[path][met]
	}
	return nil
}

// getHandlerBe Поиск и получение общих обработчиков (middleware)
func (hp *HandlerHttp) getHandlerBe(path string) (hf []HandlerFuncHttp) {
	if _, ok := hp.routeBe[path]; ok {
		return hp.routeBe[path]
	}
	return nil
}

// getHandlerAf Поиск и получение общих обработчиков (middleware)
func (hp *HandlerHttp) getHandlerAf(path string) (hf []HandlerFuncHttp) {
	if _, ok := hp.routeAf[path]; ok {
		return hp.routeAf[path]
	}
	return nil
}

// ServeHTTP Точка входа запроса (в приложение).
func (hp *HandlerHttp) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var (
		err       error
		handler   = hp.getHandler(r.URL.Path, r.Method)
		handlerBe = hp.getHandlerBe(r.URL.Path)
		handlerAf = hp.getHandlerAf(r.URL.Path)
		rw        = core.NewRW(w, r)
	)
	defer r.Body.Close()

	// Обработчики не найдены. Статика.
	if handler == nil && handlerBe == nil && handlerAf == nil {
		if err = rw.ResponseStatic(core.Config.DirWww + r.URL.Path); err != nil {
			lg.Error(err)
		}
		return
	}

	// Контекст
	ctx, cancel := context.WithTimeout(context.Background(), hp.server.WriteTimeout)
	defer cancel()

	// Обработчики
	for i, _ := range handlerBe {
		ctx, rw = handlerBe[i](ctx, rw)
	}
	if handler != nil {
		ctx, rw = handler(ctx, rw)
	}
	for i, _ := range handlerAf {
		ctx, rw = handlerAf[i](ctx, rw)
	}
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

// func (hp *HandlerHttp) String() {
// 	fmt.Println("middleware")
// 	for i, v := range hp.routeMW {
// 		fmt.Printf("%s, %#v\n", i, v)
// 	}
// 	fmt.Println("handler")
// 	for i, v := range hp.route {
// 		fmt.Printf("%s, %#v\n", i, v)
// 	}
// }
