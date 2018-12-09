package core

import (
	"net/http"
	"net/url"

	"gopkg.in/sungora/app.v1/lg"
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
	RW *rw
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
	// initialization session
	if 0 < Config.Main.SessionTimeout {
		token := self.RW.CookieGet(tool.ServiceName)
		if token == "" {
			token = tool.NewPass(10)
			self.RW.CookieSet(tool.ServiceName, token)
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

// Контроллер для реализации вывода html страниц
type ControllerHtml struct {
	Controller
	Session       *Session
	Variables     map[string]interface{} // Variable (по умолчанию пустой)
	Functions     map[string]interface{} // html/template.FuncMap (по умолчанию пустой)
	TplController string
	TplLayout     string
}

func (self *ControllerHtml) Init(w http.ResponseWriter, r *http.Request) {
	self.RW = newRW(r, w)
	self.Functions = make(map[string]interface{})
	self.Variables = make(map[string]interface{})
	// get parametrs // post "application/x-www-form-urlencoded":
	r.ParseForm()
	self.RW.RequestParams, _ = url.ParseQuery(r.URL.Query().Encode())
	for i, v := range r.Form {
		self.RW.RequestParams[i] = v
	}
	// initialization session
	if 0 < Config.Main.SessionTimeout {
		token := self.RW.CookieGet(tool.ServiceName)
		if token == "" {
			token = tool.NewPass(10)
			self.RW.CookieSet(tool.ServiceName, token)
		}
		self.Session = GetSession(token)
	}
	//
	self.TplLayout = tool.DirTpl + "/layout/new.html"
	self.TplController = tool.DirTpl + "/controllers"
}

func (self *ControllerHtml) Response() {
	if self.RW.isResponse {
		return
	}
	// шаблон контроллера
	data, err := tool.HtmlCompilation(self.TplController, self.Functions, self.Variables)
	if err != nil {
		lg.Error(err.Error())
		self.RW.ResponseHtml("internal server error", 500)
		return
	}
	// шаблон макета
	Variables := make(map[string]interface{})
	Variables["Content"] = data
	data, err = tool.HtmlCompilation(self.TplLayout, self.Functions, Variables)
	if err != nil {
		lg.Error(err.Error())
		self.RW.ResponseHtml("internal server error", 500)
		return
	}
	self.RW.ResponseHtml(data, 200)
}
