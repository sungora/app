package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"

	"github.com/sungora/app/keys"
)

// Структура для работы с входящим запросом
type Incoming struct {
	request  *http.Request
	response http.ResponseWriter
	Param    Param
}

// NewIn Функционал по работе с входящим запросом
func NewIn(w http.ResponseWriter, r *http.Request) *Incoming {
	var rw = &Incoming{
		request:  r,
		response: w,
		Param: Param{
			request: r,
		},
	}
	return rw
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

// GetBodyJson декодирование полученного тела запроса в формате json в объект
func (rw *Incoming) GetBodyJson(object interface{}) (err error) {
	var body []byte
	if body, err = ioutil.ReadAll(rw.request.Body); err != nil {
		return
	}
	if 0 == len(body) {
		return errEmptyData
	}
	return json.Unmarshal(body, object)
}

// Error ответ на запрос с ошибкой
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Data ответ на запрос с данными
type Data struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// JsonError отрицательный ответ с ошибкой в формате json (структурированный)
func (rw *Incoming) JsonError(code int, message string, status ...int) {
	res := new(Error)
	res.Code = code
	res.Message = message
	if len(status) == 0 {
		rw.Json(res, http.StatusBadRequest)
	} else {
		rw.Json(res, status[0])
	}
}

// Json ответ в формате json
func (rw *Incoming) Json(object interface{}, status ...int) {
	data, err := json.Marshal(object)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}
	// headers
	rw.generalHeaderSet("application/json; charset=utf-8", len(data))
	// Статус ответа
	if len(status) == 0 {
		rw.response.WriteHeader(http.StatusOK)
		rw.response.Header().Set(keys.Hand.Status, strconv.Itoa(http.StatusOK))
	} else {
		rw.response.WriteHeader(status[0])
		rw.response.Header().Set(keys.Hand.Status, strconv.Itoa(status[0]))
	}
	if ctx := rw.request.Context(); ctx != nil {
		if ctxChi := chi.RouteContext(ctx); ctxChi != nil {
			rw.response.Header().Set(keys.Hand.RoutePattern, strings.TrimRight(ctxChi.RoutePattern(), "/"))
		}
	}
	// Тело документа
	_, err = rw.response.Write(data)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
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
	_, _ = rw.response.Write(data)
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
		if fi, err = os.Stat(path); err != nil {
			rw.Html("<H1>Not Found</H1>", http.StatusNotFound)
			return
		}
	}
	// content
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		rw.Html("<H1>Internal Server Error</H1>", http.StatusInternalServerError)
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
