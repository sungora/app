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
	request       *http.Request
	response      http.ResponseWriter
	RequestParams map[string][]string
}

// NewRW Функционал по непосредственной работе с запросом и ответом
func NewRW(w http.ResponseWriter, r *http.Request) *RW {
	var rw = &RW{
		request:  r,
		response: w,
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
	sessionID, err := rw.request.Cookie(name)
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
	cookie.Domain = rw.request.URL.Host
	cookie.Path = `/`
	if 0 < len(t) {
		cookie.Expires = t[0]
	}
	http.SetCookie(rw.response, cookie)
}

// CookieRem Удаление куков.
func (rw *RW) CookieRem(name string) {
	var cookie = new(http.Cookie)
	cookie.Name = name
	cookie.Domain = rw.request.URL.Host
	cookie.Path = `/`
	cookie.Expires = time.Now().In(Config.TimeLocation)
	http.SetCookie(rw.response, cookie)
}

func (rw *RW) RequestBodyDecodeJson(object interface{}) (err error) {
	var body []byte
	if body, err = ioutil.ReadAll(rw.request.Body); err != nil {
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
	// запрет кеширования
	rw.response.Header().Set("Cache-Control", "no-cache, must-revalidate")
	rw.response.Header().Set("Pragma", "no-cache")
	rw.response.Header().Set("Date", t.Format(time.RFC3339))
	rw.response.Header().Set("Last-Modified", t.Format(time.RFC3339))
	// размер и тип контента
	rw.response.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.response.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	// Статус ответа
	rw.response.WriteHeader(status)
	// Тело документа
	_, err = rw.response.Write(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func (rw *RW) ResponseHtml(con string, status int) {
	data := []byte(con)
	//
	t := time.Now().In(Config.TimeLocation)
	// запрет кеширования
	rw.response.Header().Set("Cache-Control", "no-cache, must-revalidate")
	rw.response.Header().Set("Pragma", "no-cache")
	rw.response.Header().Set("Date", t.Format(time.RFC3339))
	rw.response.Header().Set("Last-Modified", t.Format(time.RFC3339))
	// размер и тип контента
	rw.response.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.response.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	// Статус ответа
	rw.response.WriteHeader(status)
	// Тело документа
	rw.response.Write(data)
}

func (rw *RW) ResponseStatic(path string) (err error) {
	var fi os.FileInfo
	if fi, err = os.Stat(path); err != nil {
		rw.ResponseHtml("<H1>Internal Server Error</H1>", http.StatusInternalServerError)
		return
	}
	if fi.IsDir() == true {
		if rw.request.URL.Path != "/" {
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
	rw.response.Header().Set("Cache-Control", "no-cache, must-revalidate")
	rw.response.Header().Set("Pragma", "no-cache")
	rw.response.Header().Set("Date", t.Format(time.RFC3339))
	rw.response.Header().Set("Last-Modified", t.Format(time.RFC3339))
	// размер и тип контента
	rw.response.Header().Set("Content-Type", typ)
	rw.response.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	// Аттач если документ не картинка и не текстововой
	if strings.LastIndex(typ, `image`) == -1 && strings.LastIndex(typ, `text`) == -1 {
		rw.response.Header().Set("Content-Disposition", "attachment; filename = "+filepath.Base(path))
	}
	// Статус ответа
	rw.response.WriteHeader(http.StatusOK)
	// Тело документа
	_, err = rw.response.Write(data)
	return
}
