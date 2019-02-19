package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RW struct {
	Request       *http.Request
	RequestParams map[string][]string
	Response      http.ResponseWriter
	Functions     map[string]interface{} // html/template.FuncMap (по умолчанию пустой)
	Variables     map[string]interface{} // Variable (по умолчанию пустой)
	TplLayout     string
	TplController string
	isResponse    bool
	Status        int
}

// NewRW Функционал по непосредственной работе с запросом и ответом
func NewRW(w http.ResponseWriter, r *http.Request) *RW {
	var rw = &RW{
		Request:       r,
		Response:      w,
		Functions:     make(map[string]interface{}),
		Variables:     make(map[string]interface{}),
		TplLayout:     Config.DirWww + "/layout",
		TplController: Config.DirWww + "/controllers",
		Status:        http.StatusOK,
	}
	// request parameter "application/x-www-form-urlencoded"
	rw.RequestParams, _ = url.ParseQuery(r.URL.Query().Encode())
	if err := r.ParseForm(); err != nil {
		return rw
	}
	for i, v := range r.Form {
		rw.RequestParams[i] = v
	}
	return rw
}

// CookieGet Получение куки.
func (rw *RW) CookieGet(name string) (c string, err error) {
	sessionID, err := rw.Request.Cookie(name)
	if err == http.ErrNoCookie {
		return "", nil
	} else if err != nil {
		return
	}
	return sessionID.Value, nil
}

// CookieSet Установка куки. Если время не указано кука сессионная (пока открыт браузер).
func (rw *RW) CookieSet(name, value string, t ...time.Time) {
	var cookie = new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Domain = rw.Request.URL.Host
	cookie.Path = `/`
	if 0 < len(t) {
		cookie.Expires = t[0]
	}
	http.SetCookie(rw.Response, cookie)
}

// CookieRem Удаление куков.
func (rw *RW) CookieRem(name string) {
	var cookie = new(http.Cookie)
	cookie.Name = name
	cookie.Domain = rw.Request.URL.Host
	cookie.Path = `/`
	cookie.Expires = time.Now().In(Config.TimeLocation)
	http.SetCookie(rw.Response, cookie)
}

func (rw *RW) RequestBodyDecodeJson(object interface{}) (err error) {
	var body []byte
	if body, err = ioutil.ReadAll(rw.Request.Body); err != nil {
		return
	}
	if 0 == len(body) {
		return errors.New("Запрос пустой, данные отсутствуют")
	}
	return json.Unmarshal(body, object)
}

type dataApi struct {
	Code    int
	Message string
	Error   bool
	Data    interface{} `json:"Data,omitempty"`
}

func (rw *RW) ResponseJsonApi200(object interface{}, code int, message string) {
	res := new(dataApi)
	res.Code = code
	res.Message = message
	res.Error = false
	res.Data = object
	rw.ResponseJson(res, http.StatusOK)
}

func (rw *RW) ResponseJsonApi409(object interface{}, code int, message string) {
	res := new(dataApi)
	res.Code = code
	res.Message = message
	res.Error = true
	res.Data = object
	rw.ResponseJson(res, http.StatusConflict)
}

func (rw *RW) ResponseJson(object interface{}, status int) {
	data, err := json.Marshal(object)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	//
	t := time.Now().In(Config.TimeLocation)
	rw.isResponse = true
	// запрет кеширования
	rw.Response.Header().Set("Cache-Control", "no-cache, must-revalidate")
	rw.Response.Header().Set("Pragma", "no-cache")
	rw.Response.Header().Set("Date", t.Format(time.RFC3339))
	rw.Response.Header().Set("Last-Modified", t.Format(time.RFC3339))
	// размер и тип контента
	rw.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	// Статус ответа
	rw.Response.WriteHeader(status)
	// Тело документа
	_, err = rw.Response.Write(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func (rw *RW) ResponseHtml(con string, status int) {
	data := []byte(con)
	//
	t := time.Now().In(Config.TimeLocation)
	rw.isResponse = true
	// запрет кеширования
	rw.Response.Header().Set("Cache-Control", "no-cache, must-revalidate")
	rw.Response.Header().Set("Pragma", "no-cache")
	rw.Response.Header().Set("Date", t.Format(time.RFC3339))
	rw.Response.Header().Set("Last-Modified", t.Format(time.RFC3339))
	// размер и тип контента
	rw.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	// Статус ответа
	rw.Response.WriteHeader(status)
	// Тело документа
	rw.Response.Write(data)
}

func (rw *RW) ResponseStatic(path string) (err error) {
	var fi os.FileInfo
	if fi, err = os.Stat(path); err != nil {
		rw.ResponseHtml("<H1>Internal Server Error</H1>", http.StatusInternalServerError)
		return
	}
	if fi.IsDir() == true {
		if rw.Request.URL.Path != "/" {
			path += string(os.PathSeparator)
		}
		path += "index.html"
	}
	// content
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		if fi.IsDir() == true {
			rw.ResponseHtml("<H1>Forbidden</H1>", http.StatusForbidden)
		} else if fi.Mode().IsRegular() == true {
			rw.ResponseHtml("<H1>Internal Server Error</H1>", http.StatusInternalServerError)
		} else {
			rw.ResponseHtml("<H1>Not Found</H1>", http.StatusNotFound)
		}
		return
	}
	// type
	var typ = `application/octet-stream`
	l := strings.Split(path, ".")
	fileExt := `.` + l[len(l)-1]
	if mimeType := mime.TypeByExtension(fileExt); mimeType != `` {
		typ = mimeType
	}
	// headers
	t := time.Now().In(Config.TimeLocation)
	// запрет кеширования
	rw.Response.Header().Set("Cache-Control", "no-cache, must-revalidate")
	rw.Response.Header().Set("Pragma", "no-cache")
	rw.Response.Header().Set("Date", t.Format(time.RFC3339))
	rw.Response.Header().Set("Last-Modified", t.Format(time.RFC3339))
	// размер и тип контента
	rw.Response.Header().Set("Content-Type", typ)
	rw.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	// Аттач если документ не картинка и не текстововой
	if strings.LastIndex(typ, `image`) == -1 && strings.LastIndex(typ, `text`) == -1 {
		rw.Response.Header().Set("Content-Disposition", "attachment; filename = "+filepath.Base(path))
	}
	// Статус ответа
	rw.Response.WriteHeader(http.StatusOK)
	// Тело документа
	_, err = rw.Response.Write(data)
	return
}
