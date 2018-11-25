package core

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/sungora/app.v1/conf"
	"gopkg.in/sungora/app.v1/lg"
)

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
	data, err := json.Marshal(object)
	if err != nil {
		lg.Error(err.Error())
		return nil
	}
	t := time.Now().In(conf.TimeLocation)
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

func (self *rw) ResponseHtml(con []byte, status int) (err error) {
	if status < 400 {
		lg.Info(status, self.Request.Method, self.Request.URL.Path)
	} else {
		lg.Error(status, self.Request.Method, self.Request.URL.Path)
	}
	t := time.Now().In(conf.TimeLocation)
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
	self.isResponse = true
	return
}

func (self *rw) ResponseImg(filePath string) (err error) {
	self.isResponse = true
	return
}

func (self *rw) ResponseFile(filePath string) (err error) {
	self.isResponse = true
	return
}
