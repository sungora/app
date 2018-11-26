// Стандартный вебсервер работающий по протоколу http
package service

import (
	"errors"
	"fmt"
	"gopkg.in/sungora/app.v1/lg"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/sungora/app.v1/core"
	"gopkg.in/sungora/app.v1/tool"
)

// newHTTP создание и запуск сервера
func newHttp() (store net.Listener, err error) {
	Server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", core.Config.Host, core.Config.Port),
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
		path    = tool.DirStatic + self.r.URL.Path
	)

	// static
	// if err = self.ResponseStatic(path, 200); err == nil {
	// 	return
	// }

	// search controller (404)
	if control, err = core.GetRoute(r.URL.Path); err != nil {
		path = tool.DirStatic + "/404.html"
		self.ResponseStatic(path, 404)
		return
	}

	// init controller
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
		path = tool.DirStatic + "/404.html"
		self.ResponseStatic(path, 404)
		return
	}

	// response controller
	control.Response()
}

func (self *httpHandler) ResponseStatic(path string, status int) (err error) {
	var fi os.FileInfo
	if fi, err = os.Lstat(path); err == nil {
		if fi.IsDir() == true {
			if self.r.URL.Path != "/" {
				path += string(os.PathSeparator)
			}
			path += "index.html"
		}
		// content
		var data []byte
		if data, err = ioutil.ReadFile(path); err != nil {
			if fi.IsDir() == true {
				return errors.New("not found: " + path)
			}
			self.w.Header().Set("Content-Type", "text/html; charset=utf-8")
			self.w.WriteHeader(500)
			self.w.Write([]byte(err.Error()))
			lg.Error(500, self.r.Method, self.r.URL.Path)
			return nil
		}
		// type
		typ := `application/octet-stream`
		l := strings.Split(path, ".")
		fileExt := `.` + l[len(l)-1]
		if mimeType := mime.TypeByExtension(fileExt); mimeType != `` {
			typ = mimeType
		}
		// headers
		t := time.Now().In(tool.TimeLocation)
		d := t.Format(time.RFC1123)
		// запрет кеширования
		self.w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		self.w.Header().Set("Pragma", "no-cache")
		self.w.Header().Set("Date", d)
		self.w.Header().Set("Last-Modified", d)
		// размер и тип контента
		self.w.Header().Set("Content-Type", typ)
		self.w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		// Аттач если документ не картинка и не текстововой
		if strings.LastIndex(typ, `image`) == -1 && strings.LastIndex(typ, `text`) == -1 {
			self.w.Header().Set("Content-Disposition", "attachment; filename = "+filepath.Base(path))
		}
		// Статус ответа
		self.w.WriteHeader(status)
		// Тело документа
		self.w.Write(data)
		//
		if status < 400 {
			lg.Info(status, self.r.Method, self.r.URL.Path)
		} else {
			lg.Error(status, self.r.Method, self.r.URL.Path)
		}
		return nil
	}
	return err
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
