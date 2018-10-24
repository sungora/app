// Стандартный вебсервер работающий по протоколу http
package app

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"time"

	"gopkg.in/sungora/app.v1/conf"
	"gopkg.in/sungora/app.v1/core"
	"gopkg.in/sungora/app.v1/lg"
)

// newHTTP создание и запуск сервера
func newWeb(c *conf.Config) (store net.Listener, err error) {
	Server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port),
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
		return nil, errors.New("http server start unknown error")
	}
	return nil, err
}

type httpHandler struct{}

// ServeHTTP Точка входа запроса (в приложение).
func (self *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var (
		control core.ControllerFace
		err     error
	)

	// search controller
	if control, err = core.GetRoute(r.URL.Path); err != nil {
		lg.Error(err.Error())
		return
	}

	// init controller
	control.Init(w, r)

	// search controller method
	objValue := reflect.ValueOf(control)
	met := objValue.MethodByName(r.Method)
	if met.IsValid() == false {
		lg.Error("not found method [%s] of control", r.Method)
		return
	}

	// пример передачи параметров в метод
	var in = make([]reflect.Value, 0)
	var params []interface{}
	for i := range params {
		in = append(in, reflect.ValueOf(params[i]))
	}

	// execute method of controller
	out := met.Call(in)
	if nil != out[0].Interface() {
		lg.Error(out[0].Interface().(error))
	}

	lg.Info(106, 200, r.Method, r.URL.Path)

	// ответ
	control.Render()

}
