package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
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

// l := strings.Split(r.Header.Get("Content-Type"), ";")
// self.ContentType = l[0]
// switch l[0] {
// case "text/plain":
// case "application/x-www-form-urlencoded":
// case "multipart/form-data":
// case "application/json":

// ////

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

func RequestGetParams(postData map[string]interface{}) string {
	u := new(url.URL)
	q := u.Query()
	for k, v := range postData {
		switch v1 := v.(type) {
		case uint64:
			q.Add(k, strconv.FormatUint(v1, 10))
		case int64:
			q.Add(k, strconv.FormatInt(v1, 10))
		case int:
			q.Add(k, strconv.Itoa(v1))
		case float64:
			q.Add(k, strconv.FormatFloat(v1, 'f', -1, 64))
		case bool:
			q.Add(k, strconv.FormatBool(v1))
		case string:
			q.Add(k, v1)
		}
	}
	// query = strings.TrimLeft(query, "&")
	return q.Encode()
}

func NewRequestGET(url string, requestBody, responseBody interface{}) (err error) {
	return NewRequest(url, "GET", requestBody, responseBody)
}

func NewRequestPOST(url string, requestBody, responseBody interface{}) (err error) {
	return NewRequest(url, "POST", requestBody, responseBody)
}

func NewRequestPUT(url string, requestBody, responseBody interface{}) (err error) {
	return NewRequest(url, "PUT", requestBody, responseBody)
}

func NewRequestDELETE(url string, requestBody, responseBody interface{}) (err error) {
	return NewRequest(url, "DELETE", requestBody, responseBody)
}

func NewRequestOPTIONS(url string, requestBody, responseBody interface{}) (err error) {
	return NewRequest(url, "OPTIONS", requestBody, responseBody)
}

func NewRequest(url, method string, requestBody, responseBody interface{}) (err error) {
	var request *http.Request
	var response *http.Response
	body := new(bytes.Buffer)
	//
	if method == "POST" || method == "PUT" {
		var data []byte
		if data, err = json.Marshal(requestBody); err != nil {
			return
		}
		if _, err = body.Write(data); err != nil {
			return
		}
	} else if p, ok := requestBody.(map[string]interface{}); ok {
		url += "?" + RequestGetParams(p)
	}
	//
	if request, err = http.NewRequest(method, url, body); err == nil {
		request.Header.Set("Content-Type", "application/json")
		c := http.Client{}
		if response, err = c.Do(request); err == nil {
			defer response.Body.Close()
			bodyResponse, _ := ioutil.ReadAll(response.Body)
			json.Unmarshal(bodyResponse, responseBody)
		}
	}
	return err
}

// ////

type content struct {
	Code    int
	Message string
	Error   bool
	Data    interface{} `json:"Data,omitempty"`
}

func (self *rw) ResponseJsonApi200(object interface{}, code int, message string) error {
	res := new(content)
	res.Code = code
	res.Message = message
	res.Error = false
	res.Data = object
	return self.ResponseJson(res, 200)
}

func (self *rw) ResponseJsonApi409(object interface{}, code int, message string) error {
	res := new(content)
	res.Code = code
	res.Message = message
	res.Error = true
	res.Data = object
	lg.Error(message)
	return self.ResponseJson(res, 409)
}

func (self *rw) ResponseJson(object interface{}, status int) (err error) {
	if status < 400 {
		lg.Info(status, self.Request.Method, self.Request.URL.Path)
	} else {
		lg.Error(status, self.Request.Method, self.Request.URL.Path)
	}
	//
	data, err := json.Marshal(object)
	if err != nil {
		lg.Error(err.Error())
		return nil
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

// func (self *rw) ResponseHtml(con []byte, status int) (err error) {
// 	if status < 400 {
// 		lg.Info(status, self.Request.Method, self.Request.URL.Path)
// 	} else {
// 		lg.Error(status, self.Request.Method, self.Request.URL.Path)
// 	}
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

func (self *rw) ResponseHtml(pathTpl string, functions, variables map[string]interface{}, status int) (err error) {
	if status < 400 {
		lg.Info(status, self.Request.Method, self.Request.URL.Path)
	} else {
		lg.Error(status, self.Request.Method, self.Request.URL.Path)
	}
	//
	var tpl *template.Template
	tpl, err = template.New(filepath.Base(pathTpl)).Funcs(functions).ParseFiles(pathTpl)
	if err != nil {
		lg.Error(err.Error())
		return nil
	}
	var ret bytes.Buffer
	if err = tpl.Execute(&ret, variables); err != nil {
		lg.Error(err.Error())
		return nil
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
	self.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	self.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(ret.Bytes())))
	// Статус ответа
	self.Response.WriteHeader(status)
	// Тело документа
	self.Response.Write(ret.Bytes())
	self.isResponse = true
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

func (self *rw) ResponseStatic(path string) (err error) {
	if _, err = os.Lstat(path); err == nil {
		// content
		var data []byte
		if data, err = ioutil.ReadFile(path); err != nil {
			self.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
			self.Response.WriteHeader(500)
			self.Response.Write([]byte(err.Error()))
			lg.Error(500, self.Request.Method, self.Request.URL.Path)
			return err
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
		return nil
	}
	self.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	self.Response.WriteHeader(500)
	self.Response.Write([]byte(err.Error()))
	lg.Error(500, self.Request.Method, self.Request.URL.Path)
	return err
}
