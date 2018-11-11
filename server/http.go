// Стандартный вебсервер работающий по протоколу http
package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"gopkg.in/sungora/app.v1/conf"
	"gopkg.in/sungora/app.v1/core"
	"gopkg.in/sungora/app.v1/lg"
)

// newHTTP создание и запуск сервера
func newHttp(c *conf.ConfigMain) (store net.Listener, err error) {
	Server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", c.Host, c.Port),
		Handler:        newHttpHandler(c),
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

type httpHandler struct {
	config *conf.ConfigMain
}

func newHttpHandler(c *conf.ConfigMain) *httpHandler {
	self := new(httpHandler)
	self.config = c
	self.config.SessionTimeout *= time.Second
	return self
}

// ServeHTTP Точка входа запроса (в приложение).
func (self *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		control core.ControllerFace
		fi      os.FileInfo
		err     error
	)

	// static
	path := conf.DirStatic + r.URL.Path
	if fi, err = os.Lstat(path); err == nil {
		var data []byte
		if fi.IsDir() == true {
			if r.URL.Path != "/" {
				path += string(os.PathSeparator)
			}
			path += "index.html"
		}
		if data, err = ioutil.ReadFile(path); err == nil {
			control := &core.Controller{}
			control.Init(w, r, self.config)
			control.RW.ResponseHtml(data, 200)
			lg.Info(200, r.Method, path)
			return
		}
	}

	// search controller (404)
	if control, err = core.GetRoute(r.URL.Path); err != nil {
		control := &core.Controller{}
		control.Init(w, r, self.config)
		control.Response(404)
		lg.Error(404, r.Method, r.URL.Path)
		return
	}

	// init controller
	control.Init(w, r, self.config)

	// init session
	if 0 < self.config.SessionTimeout {
		control.SessionStart()
	}

	// execute controller
	switch r.Method {
	case "GET":
		err = control.GET()
	case "POST":
		err = control.POST()
	case "PUT":
		err = control.PUT()
	case "DELETE":
		err = control.DELETE()
	case "OPTIONS":
		err = control.OPTIONS()
	default:
		control.Response(404)
		lg.Error(404, r.Method, r.URL.Path)
		return
	}

	// response controller
	if err != nil {
		control.Response(409)
		lg.Error(409, r.Method, r.URL.Path, err.Error())
	} else {
		control.Response(200)
		lg.Info(200, r.Method, r.URL.Path)
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
}
