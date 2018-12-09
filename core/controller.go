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

// Контроллер для реализации api запросов в формате json
type ControllerJson struct {
	RW      *rw
	Session *sessionTyp
	Data    interface{}
}

func (self *ControllerJson) Init(w http.ResponseWriter, r *http.Request) {
	self.RW = newRW(r, w)
	// request parameter "application/x-www-form-urlencoded"
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

func (self *ControllerJson) GET() {
}
func (self *ControllerJson) POST() {
}
func (self *ControllerJson) PUT() {
}
func (self *ControllerJson) DELETE() {
}
func (self *ControllerJson) OPTIONS() {
}

func (self *ControllerJson) Response() {
	if self.RW.isResponse {
		return
	}
	self.RW.ResponseJson(self.Data, self.RW.Status)
}

// Контроллер для реализации вывода html страниц
type ControllerHtml struct {
	RW            *rw
	Session       *sessionTyp
	Variables     map[string]interface{} // Variable (по умолчанию пустой)
	Functions     map[string]interface{} // html/template.FuncMap (по умолчанию пустой)
	TplController string
	TplLayout     string
}

func (self *ControllerHtml) Init(w http.ResponseWriter, r *http.Request) {
	self.RW = newRW(r, w)
	self.Functions = make(map[string]interface{})
	self.Variables = make(map[string]interface{})
	// request parameter "application/x-www-form-urlencoded"
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

func (self *ControllerHtml) GET() {
}
func (self *ControllerHtml) POST() {
}
func (self *ControllerHtml) PUT() {
}
func (self *ControllerHtml) DELETE() {
}
func (self *ControllerHtml) OPTIONS() {
}

func (self *ControllerHtml) Response() {
	if self.RW.isResponse {
		return
	}
	// шаблон контроллера
	data, err := tool.HtmlCompilation(self.TplController, self.Functions, self.Variables)
	if err != nil {
		lg.Error(err.Error())
		self.RW.ResponseHtml("", 500)
		return
	}
	// шаблон макета
	if self.TplLayout == "" {
		self.TplLayout = tool.DirTpl + "/layout/index.html"
	}
	Variables := make(map[string]interface{})
	Variables["Content"] = data
	data, err = tool.HtmlCompilation(self.TplLayout, self.Functions, Variables)
	if err != nil {
		lg.Error(err.Error())
		self.RW.ResponseHtml("", 500)
		return
	}
	self.RW.ResponseHtml(data, 200)
}
