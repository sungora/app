package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/sungora/app.v1/lg"
	"gopkg.in/sungora/app.v1/tool"
)

type rw struct {
	Request       *http.Request
	RequestParams map[string][]string
	Response      http.ResponseWriter
	isResponse    bool
	Status        int
}

func newRW(r *http.Request, w http.ResponseWriter) *rw {
	self := new(rw)
	self.Request = r
	self.Response = w
	self.Status = 200
	return self
}

// CookieGet Получение куки.
func (self *rw) CookieGet(name string) string {
	d := self.Request.Header.Get("Cookie")
	if d != "" {
		sl := strings.Split(d, ";")
		for _, v := range sl {
			sl := strings.Split(v, "=")
			sl[0] = strings.TrimSpace(sl[0])
			if sl[0] == name {
				return sl[1]
			}
		}
	}
	return ""
}

// CookieSet Установка куки. Если время не указано кука сессионная (пока открыт браузер).
func (self *rw) CookieSet(name, value string, t ...time.Time) {
	var cookie = new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Domain = self.Request.URL.Host
	cookie.Path = `/`
	if 0 < len(t) {
		cookie.Expires = t[0]
		lg.Info(101, name, value)
	} else {
		lg.Info(100, name, value)
	}
	http.SetCookie(self.Response, cookie)
}

// CookieRem Удаление куков.
func (self *rw) CookieRem(name string) {
	var cookie = new(http.Cookie)
	cookie.Name = name
	cookie.Domain = self.Request.URL.Host
	cookie.Path = `/`
	cookie.Expires = time.Now().In(tool.TimeLocation)
	http.SetCookie(self.Response, cookie)
	lg.Info(175, name)
}

func (self *rw) RequestBodyDecodeJson(object interface{}) error {
	if body, err := ioutil.ReadAll(self.Request.Body); err == nil {
		if err = json.Unmarshal(body, object); err != nil {
			return lg.Error(err.Error())
		}
	} else {
		return lg.Error(err.Error())
	}
	return nil
}

type content struct {
	Code    int
	Message string
	Error   bool
	Data    interface{} `json:"Data,omitempty"`
}

func (self *rw) ResponseJsonApi200(object interface{}, code int, message string) {
	res := new(content)
	res.Code = code
	res.Message = message
	res.Error = false
	res.Data = object
	self.ResponseJson(res, 200)
}

func (self *rw) ResponseJsonApi403(object interface{}, code int, message string) {
	res := new(content)
	res.Code = code
	res.Message = message
	res.Error = false
	res.Data = object
	self.ResponseJson(res, 403)
}

func (self *rw) ResponseJsonApi409(object interface{}, code int, message string) {
	res := new(content)
	res.Code = code
	res.Message = message
	res.Error = true
	res.Data = object
	self.ResponseJson(res, 409)
}

func (self *rw) ResponseJson(object interface{}, status int) {
	if status < 400 {
		lg.Info(status, self.Request.Method, self.Request.URL.Path)
	} else {
		lg.Error(status, self.Request.Method, self.Request.URL.Path)
	}
	//
	data, err := json.Marshal(object)
	if err != nil {
		lg.Error(err.Error())
		return
	}
	//
	t := time.Now().In(tool.TimeLocation)
	d := t.Format(time.RFC1123)
	// запрет кеширования
	self.Response.Header().Set("Cache-Control", "no-cache, must-revalidate")
	self.Response.Header().Set("Pragma", "no-cache")
	self.Response.Header().Set("Date", d)
	self.Response.Header().Set("Last-Modified", d)
	// размер и тип контента
	self.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	self.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	// Статус ответа
	self.Response.WriteHeader(status)
	// Тело документа
	self.Response.Write(data)
	self.isResponse = true
	return
}

func (self *rw) ResponseHtml(con string, status int) {
	if status < 400 {
		lg.Info(status, self.Request.Method, self.Request.URL.Path)
	} else {
		lg.Error(status, self.Request.Method, self.Request.URL.Path)
	}
	t := time.Now().In(tool.TimeLocation)
	d := t.Format(time.RFC1123)
	// запрет кеширования
	self.Response.Header().Set("Cache-Control", "no-cache, must-revalidate")
	self.Response.Header().Set("Pragma", "no-cache")
	self.Response.Header().Set("Date", d)
	self.Response.Header().Set("Last-Modified", d)
	// размер и тип контента
	self.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	self.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(con)))
	// Статус ответа
	self.Response.WriteHeader(status)
	// Тело документа
	self.Response.Write([]byte(con))
	self.isResponse = true
	return
}

func (self *rw) ResponseStatic(path string) {
	var err error
	var fi os.FileInfo
	if fi, err = os.Lstat(path); err == nil {
		if fi.IsDir() == true {
			if self.Request.URL.Path != "/" {
				path += string(os.PathSeparator)
			}
			path += "index.html"
		}
		// content
		var data []byte
		if data, err = ioutil.ReadFile(path); err != nil {
			if fi.IsDir() == true {
				self.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
				self.Response.WriteHeader(403)
				self.Response.Write([]byte("<h1>Access forbidden</h1>"))
				lg.Error(403, self.Request.Method, self.Request.URL.Path)
				return

			} else {
				self.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
				self.Response.WriteHeader(500)
				self.Response.Write([]byte("<h1>Internal server error</h1>"))
				lg.Error(500, self.Request.Method, self.Request.URL.Path)
				return
			}
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
		self.Response.Header().Set("Cache-Control", "no-cache, must-revalidate")
		self.Response.Header().Set("Pragma", "no-cache")
		self.Response.Header().Set("Date", d)
		self.Response.Header().Set("Last-Modified", d)
		// размер и тип контента
		self.Response.Header().Set("Content-Type", typ)
		self.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		// Аттач если документ не картинка и не текстововой
		if strings.LastIndex(typ, `image`) == -1 && strings.LastIndex(typ, `text`) == -1 {
			self.Response.Header().Set("Content-Disposition", "attachment; filename = "+filepath.Base(path))
		}
		// Статус ответа
		self.Response.WriteHeader(200)
		// Тело документа
		self.Response.Write(data)
		lg.Info(200, self.Request.Method, self.Request.URL.Path)
		return
	}
	self.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	self.Response.WriteHeader(404)
	self.Response.Write([]byte("<h1>Page not found</h1>"))
	lg.Error(404, self.Request.Method, self.Request.URL.Path)
	return
}

// func (self *rw) ResponseImg(filePath string) (err error) {
// 	lg.Info(200, self.Request.Method, self.Request.URL.Path)
// 	t := time.Now().In(tool.TimeLocation)
// 	d := t.Format(time.RFC1123)
// 	// запрет кеширования
// 	self.Response.Header().Set("Cache-Control", "no-cache, must-revalidate")
// 	self.Response.Header().Set("Pragma", "no-cache")
// 	self.Response.Header().Set("Date", d)
// 	self.Response.Header().Set("Last-Modified", d)
// 	// размер и тип контента
// 	self.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	self.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(con)))
// 	// Статус ответа
// 	self.Response.WriteHeader(status)
// 	// Тело документа
// 	self.Response.Write(con)
// 	self.isResponse = true
// 	return
// }

// func (self *rw) ResponseFile(filePath string) (err error) {
// 	lg.Info(200, self.Request.Method, self.Request.URL.Path)
// 	t := time.Now().In(tool.TimeLocation)
// 	d := t.Format(time.RFC1123)
// 	// запрет кеширования
// 	self.Response.Header().Set("Cache-Control", "no-cache, must-revalidate")
// 	self.Response.Header().Set("Pragma", "no-cache")
// 	self.Response.Header().Set("Date", d)
// 	self.Response.Header().Set("Last-Modified", d)
// 	// размер и тип контента
// 	self.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	self.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(con)))
// 	// Статус ответа
// 	self.Response.WriteHeader(status)
// 	// Тело документа
// 	self.Response.Write(con)
// 	self.isResponse = true
// 	return
// }

// ////
// l := strings.Split(r.Header.Get("Content-Type"), ";")
// self.ContentType = l[0]
// switch l[0] {
// case "text/plain":
// case "application/x-www-form-urlencoded":
// case "multipart/form-data":
// case "application/json":
// ////
