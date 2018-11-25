package core

import (
	"net/http"
	"net/url"

	"gopkg.in/sungora/app.v1/tool"
)

// ContraFace is an interface to uniform all controller handler.
type ControllerFace interface {
	Init(w http.ResponseWriter, r *http.Request)
	GET()
	POST()
	PUT()
	DELETE()
	OPTIONS()
	Response()
}

// Базовый контроллер
type Controller struct {
	RW     *rw
}

func (self *Controller) Init(w http.ResponseWriter, r *http.Request) {
	self.RW = newRW(r, w)
}

func (self *Controller) GET() {
}
func (self *Controller) POST() {
}
func (self *Controller) PUT() {
}
func (self *Controller) DELETE() {
}
func (self *Controller) OPTIONS() {
}
func (self *Controller) Response() {
}

// Контроллер для реализации api запросов в формате json
type ControllerJson struct {
	Controller
	Session *Session
	Data    interface{}
}

func (self *ControllerJson) Init(w http.ResponseWriter, r *http.Request) {
	self.RW = newRW(r, w)
	// get parametrs // post "application/x-www-form-urlencoded":
	r.ParseForm()
	self.RW.RequestParams, _ = url.ParseQuery(r.URL.Query().Encode())
	for i, v := range r.Form {
		self.RW.RequestParams[i] = v
	}
	// init session
	if 0 < Config.SessionTimeout {
		token := self.RW.CookieGet(Config.Name)
		if token == "" {
			token = tool.NewPass(10)
			self.RW.CookieSet(Config.Name, token)
		}
		self.Session = GetSession(token)
	}
}

func (self *ControllerJson) Response() {
	if self.RW.isResponse {
		return
	}
	self.RW.ResponseJson(self.Data, self.RW.Status)
}
