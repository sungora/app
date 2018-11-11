package core

import (
	"gopkg.in/sungora/app.v1/conf"
	"gopkg.in/sungora/app.v1/lg"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type rw struct {
	Request        *http.Request
	RequestParams      map[string][]string
	Response       http.ResponseWriter
	responseStatus bool
}

func newRW(r *http.Request, w http.ResponseWriter) *rw {
	self := new(rw)
	self.Request = r
	self.Response = w
	// get
	self.RequestParams, _ = url.ParseQuery(r.URL.Query().Encode())
	// post "application/x-www-form-urlencoded":
	r.ParseForm()
	for i, v := range r.Form {
		self.RequestParams[i] = v
	}
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
	cookie.Expires = time.Now().In(conf.TimeLocation)
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
