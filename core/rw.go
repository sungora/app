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
)

type requestResponse struct {
	Request       *http.Request
	RequestParams map[string][]string
	Response      http.ResponseWriter
	isResponse    bool
	Status        int
}

func NewRW(r *http.Request, w http.ResponseWriter) *requestResponse {
	rw := new(requestResponse)
	rw.Request = r
	rw.Response = w
	rw.Status = http.StatusOK
	return rw
}

// CookieGet Получение куки.
func (rw *requestResponse) CookieGet(name string) (c string, err error) {
	sessionID, err := rw.Request.Cookie(name)
	if err == http.ErrNoCookie {
		return "", nil
	} else if err != nil {
		return
	}
	return sessionID.Value, nil
}

// CookieSet Установка куки. Если время не указано кука сессионная (пока открыт браузер).
func (rw *requestResponse) CookieSet(name, value string, t ...time.Time) {
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
func (rw *requestResponse) CookieRem(name string) {
	var cookie = new(http.Cookie)
	cookie.Name = name
	cookie.Domain = rw.Request.URL.Host
	cookie.Path = `/`
	cookie.Expires = time.Now().In(TimeLocation)
	http.SetCookie(rw.Response, cookie)
}

func (rw *requestResponse) RequestBodyDecodeJson(object interface{}) (err error) {
	var body []byte
	if body, err = ioutil.ReadAll(rw.Request.Body); err != nil {
		return
	}
	if 0 == len(body) {
		return errors.New("Запрос пустой, данные отсутствуют")
	}
	return json.Unmarshal(body, object)
}

type content struct {
	Code    int
	Message string
	Error   bool
	Data    interface{} `json:"Data,omitempty"`
}

func (rw *requestResponse) ResponseJsonOk(object interface{}, code int, message string) error {
	res := new(content)
	res.Code = code
	res.Message = message
	res.Error = false
	res.Data = object
	return rw.ResponseJson(res, http.StatusOK)
}

func (rw *requestResponse) ResponseJsonErr(object interface{}, code int, message string, status ...int) error {
	res := new(content)
	res.Code = code
	res.Message = message
	res.Error = true
	res.Data = object
	if 0 < len(status) {
		return rw.ResponseJson(res, status[0])
	}
	return rw.ResponseJson(res, http.StatusConflict)
}

func (rw *requestResponse) ResponseJson(object interface{}, status int) error {
	data, err := json.Marshal(object)
	if err != nil {
		return err
	}
	//
	t := time.Now().In(TimeLocation)
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
	return err
}

func (rw *requestResponse) ResponseHtml(con string, status int) {
	data := []byte(con)
	//
	t := time.Now().In(TimeLocation)
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

func (rw *requestResponse) ResponseStatic(path string) (err error) {
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
	t := time.Now().In(TimeLocation)
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
