// Стандартный вебсервер работающий по протоколу http
package service

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"gopkg.in/sungora/app.v1/core"
	"gopkg.in/sungora/app.v1/tool"
)

// newHTTP создание и запуск сервера
func newHttp() (store net.Listener, err error) {
	Server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", core.Config.Main.Host, core.Config.Main.Port),
		Handler:        new(httpHandler),
		ReadTimeout:    time.Second * time.Duration(300),
		WriteTimeout:   time.Second * time.Duration(300),
		MaxHeaderBytes: 1048576,
	}
	for i := 5; i > 0; i-- {
		store, err = net.Listen("tcp", Server.Addr)
		time.Sleep(time.Millisecond * 100)
		if err == nil {
			break
		}
	}
	if err == nil && store != nil {
		go Server.Serve(store)
		return
	} else if err == nil {
		return nil, errors.New("service start unknown error")
	}
	return nil, err
}

type httpHandler struct {
	w http.ResponseWriter
	r *http.Request
}

// ServeHTTP Точка входа запроса (в приложение).
func (self *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.w = w
	self.r = r
	var (
		err     error
		control core.ControllerFace
	)

	// search controller & static
	if control, err = core.Route.GetRoute(r.URL.Path); err != nil {
		c := new(core.Controller)
		c.Init(w, r)
		c.RW.ResponseStatic(tool.DirWww + self.r.URL.Path)
		return
	}

	// initialization controller
	control.Init(w, r)

	// execute controller
	switch r.Method {
	case "GET":
		control.GET()
	case "POST":
		control.POST()
	case "PUT":
		control.PUT()
	case "DELETE":
		control.DELETE()
	case "OPTIONS":
		control.OPTIONS()
	default:
		c := new(core.Controller)
		c.Init(w, r)
		c.RW.ResponseStatic(tool.DirTpl + "/404.html")
		return
	}

	// response controller
	control.Response()
}

// // search controller method
// objValue := reflect.ValueOf(control)
// met := objValue.MethodByName(r.Method)
// if met.IsValid() == false {
// 	control := &core.Controller{}
// 	control.Init(w, r, self.config)
// 	control.RW.ResponseHtml([]byte("page not found (m)"), 404)
// 	lg.Error("not found method [%s] of control", r.Method)
// 	return
// }
//
// // пример передачи параметров в метод
// var in = make([]reflect.Value, 0)
// var params []interface{}
// for i := range params {
// 	in = append(in, reflect.ValueOf(params[i]))
// }
//
// // execute method of controller
// out := met.Call(in)
// if nil != out[0].Interface() {
// 	lg.Error(out[0].Interface().(error))
// }
