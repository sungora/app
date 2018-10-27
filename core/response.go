package core

import (
	"encoding/json"
	"fmt"
	"time"
)

type content struct {
	Code    int
	Message string
	Error   bool
	Data    interface{} `json:"Data,omitempty"`
}

func (self *rw) JsonApi(object interface{}, code int, message string, status int) error {
	res := new(content)
	res.Code = code
	res.Message = message
	if status > 399 {
		res.Error = true
	}
	res.Data = object
	return self.Json(res, status)
}

func (self *rw) Json200(object interface{}) error {
	return self.Json(object, 200)
}

func (self *rw) Json409(object interface{}) error {
	return self.Json(object, 409)
}

func (self *rw) Json(object interface{}, status int) (err error) {
	con, err := json.Marshal(object)
	if err != nil {
		return err
	}
	var loc *time.Location
	if loc, err = time.LoadLocation(`Europe/Moscow`); err != nil {
		loc = time.UTC
	}
	t := time.Now().In(loc)
	d := t.Format(time.RFC1123)
	// запрет кеширования
	self.Response.Header().Set("Cache-Control", "no-cache, must-revalidate")
	self.Response.Header().Set("Pragma", "no-cache")
	self.Response.Header().Set("Date", d)
	self.Response.Header().Set("Last-Modified", d)
	// размер и тип контента
	self.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	self.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(con)))
	// Статус ответа
	self.Response.WriteHeader(status)
	// Тело документа
	self.Response.Write(con)
	self.responseStatus = true
	return
}

func (self *rw) Html(con []byte, status int) (err error) {
	var loc *time.Location
	if loc, err = time.LoadLocation(`Europe/Moscow`); err != nil {
		loc = time.UTC
	}
	t := time.Now().In(loc)
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
	self.Response.Write(con)
	self.responseStatus = true
	return
}

func (self *rw) Img(filePath string) (err error) {
	self.responseStatus = true
	return
}

func (self *rw) File(filePath string) (err error) {
	self.responseStatus = true
	return
}
