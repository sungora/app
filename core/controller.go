package core

import (
	"net/http"

	"gopkg.in/sungora/app.v1/conf"
	"gopkg.in/sungora/app.v1/tool"
)

// ContraFace is an interface to uniform all controller handler.
type ControllerFace interface {
	Init(w http.ResponseWriter, r *http.Request, c *conf.ConfigMain)
	SessionStart()
	GET() (err error)
	POST() (err error)
	PUT() (err error)
	DELETE() (err error)
	OPTIONS() (err error)
	Response(status int)
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
}

// SessionStart Старт сессии
func (self *Controller) SessionStart() {
	token := self.RW.GetCookie(self.Config.Name)
	if token == "" {
		token = tool.NewPass(10)
		self.RW.SetCookie(self.Config.Name, token)
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
func (self *Controller) Response(status int) {
	if self.RW.responseStatus {
		return
	}
	switch status {
	case 403:
		self.RW.Json([]byte("Access forbidden!"), 403)
	case 404:
		self.RW.Json([]byte("Page not found"), 404)
	default:
		self.RW.Json(self.Data, status)
	}
}
