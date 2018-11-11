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

func (self *rw) ResponseJsonApi(object interface{}, code int, message string, status int) error {
	res := new(content)
	res.Code = code
	res.Message = message
	if status > 399 {
		res.Error = true
	}
	res.Data = object
	return self.ResponseJson(res, status)
}

func (self *rw) ResponseJson200(object interface{}) error {
	return self.ResponseJson(object, 200)
}

func (self *rw) ResponseJson409(object interface{}) error {
	return self.ResponseJson(object, 409)
}

func (self *rw) ResponseJson(object interface{}, status int) (err error) {
	data, err := json.Marshal(object)
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
	self.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	// Статус ответа
	self.Response.WriteHeader(status)
	// Тело документа
	self.Response.Write(data)
	self.responseStatus = true
	return
}

func (self *rw) ResponseHtml(con []byte, status int) (err error) {
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

func (self *rw) ResponseImg(filePath string) (err error) {
	self.responseStatus = true
	return
}

func (self *rw) ResponseFile(filePath string) (err error) {
	self.responseStatus = true
	return
}
