package request

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

// Структура для работы с входящим запросом
type Incoming struct {
	request       *http.Request
	response      http.ResponseWriter
	requestParams map[string][]string
}

// NewIn Функционал по работе с входящим запросом
func NewIn(w http.ResponseWriter, r *http.Request) *Incoming {
	var rw = &Incoming{
		request:  r,
		response: w,
	}
	return rw
}

// GetRequestParam Получение данных запроса пришедших в формате "application/x-www-form-urlencoded".
func (rw *Incoming) GetRequestParam(name string) map[string][]string {
	if rw.requestParams != nil {
		return rw.requestParams
	}
	rw.requestParams, _ = url.ParseQuery(rw.request.URL.Query().Encode())
	if err := rw.request.ParseForm(); err != nil {
		return rw.requestParams
	}
	for i, v := range rw.request.Form {
		rw.requestParams[i] = v
	}
	return rw.requestParams
}

// CookieGet Получение куки.
func (rw *Incoming) CookieGet(name string) (c string, err error) {
	sessionID, err := rw.request.Cookie(name)
	if err == http.ErrNoCookie {
		return "", nil
	} else if err != nil {
		return
	}
	return sessionID.Value, nil
}

// CookieSet Установка куки. Если время не указано кука сессионная (пока открыт браузер).
func (rw *Incoming) CookieSet(name, value string, t ...time.Time) {
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
func (rw *Incoming) CookieRem(name string) {
	var cookie = new(http.Cookie)
	cookie.Name = name
	cookie.Domain = rw.request.URL.Host
	cookie.Path = `/`
	cookie.Expires = time.Now()
	http.SetCookie(rw.response, cookie)
}

var errEmptyData = errors.New("Запрос пустой, данные отсутствуют")

// BodyDecodeJson декодирование полученного тела запроса в формате json в объект
func (rw *Incoming) BodyDecodeJson(object interface{}) (err error) {
	var body []byte
	if body, err = ioutil.ReadAll(rw.request.Body); err != nil {
		return
	}
	if 0 == len(body) {
		return errEmptyData
	}
	return json.Unmarshal(body, object)
}

// обертка api ответа в формате json
type JsonApi struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   bool        `json:"error"`
	Data    interface{} `json:"data,omitempty"`
}

// JsonApi200 положительный ответ api в формате json
func (rw *Incoming) JsonApi200(object interface{}, code int, message string) {
	res := new(JsonApi)
	res.Code = code
	res.Message = message
	res.Error = false
	res.Data = object
	rw.Json(res, http.StatusOK)
}

// JsonApi409 отрицательный ответ api в формате json
func (rw *Incoming) JsonApi409(object interface{}, code int, message string) {
	res := new(JsonApi)
	res.Code = code
	res.Message = message
	res.Error = true
	res.Data = object
	rw.Json(res, http.StatusConflict)
}

// Json ответ в формате json
func (rw *Incoming) Json(object interface{}, status ...int) {
	data, err := json.Marshal(object)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	// headers
	rw.generalHeaderSet("application/json; charset=utf-8", len(data))
	// Статус ответа
	if len(status) == 0 {
		rw.response.WriteHeader(http.StatusOK)
	} else {
		rw.response.WriteHeader(status[0])
	}
	// Тело документа
	_, err = rw.response.Write(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

// Html ответ в html формате
func (rw *Incoming) Html(con string, status ...int) {
	data := []byte(con)
	// headers
	rw.generalHeaderSet("text/html; charset=utf-8", len(data))
	// Статус ответа
	if len(status) == 0 {
		rw.response.WriteHeader(http.StatusOK)
	} else {
		rw.response.WriteHeader(status[0])
	}
	// Тело документа
	rw.response.Write(data)
}

// Static ответ - отдача статических данных
func (rw *Incoming) Static(path string) (err error) {
	var fi os.FileInfo
	if fi, err = os.Stat(path); err != nil {
		rw.Html("<H1>Not Found</H1>", http.StatusNotFound)
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
			rw.Html("<H1>Forbidden</H1>", http.StatusForbidden)
		} else if fi.Mode().IsRegular() == true {
			rw.Html("<H1>Internal Server Error</H1>", http.StatusInternalServerError)
		} else {
			rw.Html("<H1>Not Found</H1>", http.StatusNotFound)
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
	rw.generalHeaderSet(typ, len(data))
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

// generalHeaderSet общие заголовки любого ответа
func (rw *Incoming) generalHeaderSet(contentTyp string, l int) {
	t := time.Now()
	// запрет кеширования
	rw.response.Header().Set("Cache-Control", "no-cache, must-revalidate")
	rw.response.Header().Set("Pragma", "no-cache")
	rw.response.Header().Set("Date", t.Format(time.RFC3339))
	rw.response.Header().Set("Last-Modified", t.Format(time.RFC3339))
	// размер и тип контента
	rw.response.Header().Set("Content-Type", contentTyp)
	rw.response.Header().Set("Content-Length", fmt.Sprintf("%d", l))
}
