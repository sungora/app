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
		// TODO доработать тип отдаваемого документа
		if data, err = ioutil.ReadFile(path); err == nil {
			control := &core.Controller{}
			control.Init(w, r, self.config)
			control.RW.ResponseHtml(data, 200)
			return
		}
	}

	// search controller (404)
	if control, err = core.GetRoute(r.URL.Path); err != nil {
		control := &core.Controller{}
		control.Init(w, r, self.config)
		control.RW.Status = 404
		control.Response()
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
		control := &core.Controller{}
		control.Init(w, r, self.config)
		control.RW.Status = 404
		control.Response()
		return
	}

	// response controller
	control.Response()

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
