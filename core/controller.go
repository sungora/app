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

func (c *ControllerJson) Init(w http.ResponseWriter, r *http.Request) {
	c.RW = newRW(r, w)
	// request parameter "application/x-www-form-urlencoded"
	r.ParseForm()
	c.RW.RequestParams, _ = url.ParseQuery(r.URL.Query().Encode())
	for i, v := range r.Form {
		c.RW.RequestParams[i] = v
	}
	// initialization session
	if 0 < Config.Main.SessionTimeout {
		token := c.RW.CookieGet(tool.ServiceName)
		if token == "" {
			token = tool.NewRandomString(10)
			c.RW.CookieSet(tool.ServiceName, token)
		}
		c.Session = GetSession(token)
	}
}

func (c *ControllerJson) GET() {
}
func (c *ControllerJson) POST() {
}
func (c *ControllerJson) PUT() {
}
func (c *ControllerJson) DELETE() {
}
func (c *ControllerJson) OPTIONS() {
}

func (c *ControllerJson) Response() {
	if c.RW.isResponse {
		return
	}
	c.RW.ResponseJson(c.Data, c.RW.Status)
}

// Контроллер для реализации html страниц
type ControllerHtml struct {
	RW            *rw
	Session       *sessionTyp
	Variables     map[string]interface{} // Variable (по умолчанию пустой)
	Functions     map[string]interface{} // html/template.FuncMap (по умолчанию пустой)
	TplController string
	TplLayout     string
}

func (c *ControllerHtml) Init(w http.ResponseWriter, r *http.Request) {
	c.RW = newRW(r, w)
	// request parameter "application/x-www-form-urlencoded"
	r.ParseForm()
	c.RW.RequestParams, _ = url.ParseQuery(r.URL.Query().Encode())
	for i, v := range r.Form {
		c.RW.RequestParams[i] = v
	}
	// initialization session
	if 0 < Config.Main.SessionTimeout {
		token := c.RW.CookieGet(tool.ServiceName)
		if token == "" {
			token = tool.NewRandomString(10)
			c.RW.CookieSet(tool.ServiceName, token)
		}
		c.Session = GetSession(token)
	}
	//
	c.Functions = make(map[string]interface{})
	c.Variables = make(map[string]interface{})
	c.TplController = tool.DirTpl + "/controllers"
}

func (c *ControllerHtml) GET() {
}
func (c *ControllerHtml) POST() {
}
func (c *ControllerHtml) PUT() {
}
func (c *ControllerHtml) DELETE() {
}
func (c *ControllerHtml) OPTIONS() {
}

func (c *ControllerHtml) Response() {
	if c.RW.isResponse {
		return
	}
	// шаблон контроллера
	data, err := tool.HtmlCompilation(c.TplController, c.Functions, c.Variables)
	if err != nil {
		lg.Error(err.Error())
		c.RW.ResponseHtml("", 500)
		return
	}
	// шаблон макета
	if c.TplLayout == "" {
		c.TplLayout = tool.DirTpl + "/layout/index.html"
	}
	Variables := make(map[string]interface{})
	Variables["Content"] = data
	data, err = tool.HtmlCompilation(c.TplLayout, c.Functions, Variables)
	if err != nil {
		lg.Error(err.Error())
		c.RW.ResponseHtml("", 500)
		return
	}
	c.RW.ResponseHtml(data, 200)
}
