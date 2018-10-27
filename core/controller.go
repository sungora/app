package core

import (
	"gopkg.in/sungora/app.v1/conf"
	"net/http"
)

// ContraFace is an interface to uniform all controller handler.
type ControllerFace interface {
	Init(w http.ResponseWriter, r *http.Request, c *conf.ConfigMain)
	GET() (err error)
	POST() (err error)
	PUT() (err error)
	DELETE() (err error)
	OPTIONS() (err error)
	Render()
}

type Controller struct {
	Config  *conf.ConfigMain
	Session *Session
	RW      *rw
	Data    interface{}
}

func (self *Controller) Init(w http.ResponseWriter, r *http.Request, c *conf.ConfigMain) {
	self.Config = c
	self.RW = newRW(r, w)
	// сессия
	token := self.RW.GetCookie(c.Name)
	if token == "" {
		token = CreatePassword()
		self.RW.SetCookie(c.Name, token)
	}
	self.Session = GetSession(token)
}
func (self *Controller) GET() (err error) {
	return
}
func (self *Controller) POST() (err error) {
	return
}
func (self *Controller) PUT() (err error) {
	return
}
func (self *Controller) DELETE() (err error) {
	return
}
func (self *Controller) OPTIONS() (err error) {
	return
}
func (self *Controller) Render() {
	if self.RW.responseStatus {
		return
	}
	self.RW.Json200(self.Data)
}
