package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sungora/app/v2/lg"
	"github.com/sungora/app/v2/tool"
)

type rw struct {
	Request       *http.Request
	RequestParams map[string][]string
	Response      http.ResponseWriter
	isResponse    bool
	Status        int
}

func newRW(r *http.Request, w http.ResponseWriter) *rw {
	io := new(rw)
	io.Request = r
	io.Response = w
	io.Status = http.StatusOK
	return io
}

// CookieGet Получение куки.
func (io *rw) CookieGet(name string) string {

	sessionID, err := io.Request.Cookie(name)
	if err == http.ErrNoCookie {
		return ""
	} else if err != nil {
		lg.Error(err.Error())
		return ""
	}
	return sessionID.Value

	// d := io.Request.Header.Get("Cookie")
	// if d != "" {
	// 	sl := strings.Split(d, ";")
	// 	for _, v := range sl {
	// 		sl := strings.Split(v, "=")
	// 		sl[0] = strings.TrimSpace(sl[0])
	// 		if sl[0] == name {
	// 			return sl[1]
	// 		}
	// 	}
	// }
	// return ""
}

// CookieSet Установка куки. Если время не указано кука сессионная (пока открыт браузер).
func (io *rw) CookieSet(name, value string, t ...time.Time) {
	var cookie = new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Domain = io.Request.URL.Host
	cookie.Path = `/`
	if 0 < len(t) {
		cookie.Expires = t[0]
		lg.Info(101, name, value)
	} else {
		lg.Info(100, name, value)
	}
	http.SetCookie(io.Response, cookie)
}

// CookieRem Удаление куков.
func (io *rw) CookieRem(name string) {
	var cookie = new(http.Cookie)
	cookie.Name = name
	cookie.Domain = io.Request.URL.Host
	cookie.Path = `/`
	cookie.Expires = time.Now().In(tool.TimeLocation)
	http.SetCookie(io.Response, cookie)
	lg.Info(175, name)
}

func (io *rw) RequestBodyDecodeJson(object interface{}) error {
	if body, err := ioutil.ReadAll(io.Request.Body); err == nil {
		if 0 == len(body) {
			return errors.New("Запрос пустой, данные отсутствуют")
		}
		if err = json.Unmarshal(body, object); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

type content struct {
	Code    int
	Message string
	Error   bool
	Data    interface{} `json:"Data,omitempty"`
}

func (io *rw) ResponseJsonApi200(object interface{}, code int, message string) {
	res := new(content)
	res.Code = code
	res.Message = message
	res.Error = false
	res.Data = object
	io.ResponseJson(res, http.StatusOK)
}

func (io *rw) ResponseJsonApi403(object interface{}, code int, message string) {
	res := new(content)
	res.Code = code
	res.Message = message
	res.Error = false
	res.Data = object
	io.ResponseJson(res, http.StatusForbidden)
}

func (io *rw) ResponseJsonApi409(object interface{}, code int, message string) {
	res := new(content)
	res.Code = code
	res.Message = message
	res.Error = true
	res.Data = object
	io.ResponseJson(res, http.StatusConflict)
}

func (io *rw) ResponseJson(object interface{}, status int) {
	data, err := json.Marshal(object)
	if err != nil {
		lg.Error(err.Error())
		return
	}
	//
	t := time.Now().In(tool.TimeLocation)
	d := t.Format(time.RFC1123)
	// запрет кеширования
	io.Response.Header().Set("Cache-Control", "no-cache, must-revalidate")
	io.Response.Header().Set("Pragma", "no-cache")
	io.Response.Header().Set("Date", d)
	io.Response.Header().Set("Last-Modified", d)
	// размер и тип контента
	io.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	// Статус ответа
	io.Response.WriteHeader(status)
	// Тело документа
	io.Response.Write(data)
	io.isResponse = true
	return
}

func (io *rw) ResponseHtml(con string, status int) {
	if status < 400 {
		lg.Info(status, io.Request.Method, io.Request.URL.Path)
	} else {
		lg.Error(status, io.Request.Method, io.Request.URL.Path)
	}
	var err error
	var data []byte
	if con == "" {
		path := fmt.Sprintf("%s/layout/%d.html", tool.DirTpl, status)
		if data, err = ioutil.ReadFile(path); err != nil {
			data = []byte("<h1>Internal server error</h1>")
		}
	} else {
		data = []byte(con)
	}
	t := time.Now().In(tool.TimeLocation)
	d := t.Format(time.RFC1123)
	// запрет кеширования
	io.Response.Header().Set("Cache-Control", "no-cache, must-revalidate")
	io.Response.Header().Set("Pragma", "no-cache")
	io.Response.Header().Set("Date", d)
	io.Response.Header().Set("Last-Modified", d)
	// размер и тип контента
	io.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	// Статус ответа
	io.Response.WriteHeader(status)
	// Тело документа
	io.Response.Write(data)
	io.isResponse = true
	return
}

func (io *rw) ResponseStatic(path string) {
	var err error
	var fi os.FileInfo
	if fi, err = os.Lstat(path); err == nil {
		if fi.IsDir() == true {
			if io.Request.URL.Path != "/" {
				path += string(os.PathSeparator)
			}
			path += "index.html"
		}
		// content
		var data []byte
		if data, err = ioutil.ReadFile(path); err != nil {
			if fi.IsDir() == true {
				io.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
				io.Response.WriteHeader(http.StatusForbidden)
				io.Response.Write([]byte("<h1>Access forbidden</h1>"))
				lg.Error(403, io.Request.Method, io.Request.URL.Path)
				return

			} else {
				io.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
				io.Response.WriteHeader(http.StatusInternalServerError)
				io.Response.Write([]byte("<h1>Internal server error</h1>"))
				lg.Error(500, io.Request.Method, io.Request.URL.Path)
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
		io.Response.Header().Set("Cache-Control", "no-cache, must-revalidate")
		io.Response.Header().Set("Pragma", "no-cache")
		io.Response.Header().Set("Date", d)
		io.Response.Header().Set("Last-Modified", d)
		// размер и тип контента
		io.Response.Header().Set("Content-Type", typ)
		io.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		// Аттач если документ не картинка и не текстововой
		if strings.LastIndex(typ, `image`) == -1 && strings.LastIndex(typ, `text`) == -1 {
			io.Response.Header().Set("Content-Disposition", "attachment; filename = "+filepath.Base(path))
		}
		// Статус ответа
		io.Response.WriteHeader(http.StatusOK)
		// Тело документа
		io.Response.Write(data)
		lg.Info(http.StatusOK, io.Request.Method, io.Request.URL.Path)
		return
	}
	io.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.Response.WriteHeader(http.StatusNotFound)
	io.Response.Write([]byte("<h1>Page not found</h1>"))
	lg.Error(http.StatusNotFound, io.Request.Method, io.Request.URL.Path)
	return
}

// ////
// l := strings.Split(r.Header.Get("Content-Type"), ";")
// self.ContentType = l[0]
// switch l[0] {
// case "text/plain":
// case "application/x-www-form-urlencoded":
// case "multipart/form-data":
// case "application/json":
// ////
