package core

import (
	"net/http"
)

// ContraFace is an interface to uniform all controller handler.
type ControllerFace interface {
	Init(w http.ResponseWriter, r *http.Request)
	GET() (err error)
	POST() (err error)
	PUT() (err error)
	DELETE() (err error)
	OPTIONS() (err error)
	Render()
}

type Controller struct {
	Request  *request
	Response *response
	DataHtml map[string]interface{}
	DataJson interface{}
}

func (self *Controller) Init(w http.ResponseWriter, r *http.Request) {
	self.Request = newrequest(r)
	self.Response = newresponse(w)
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
	if self.Response.responseStatus {
		return
	}
	if self.Response.ResponseType == "json" {
		self.Response.Json200(self.DataJson)
	} else if self.Response.ResponseType == "html" {
		self.Response.Html([]byte("<H1>Html Output</H1>"), 200)
	} else {
		self.Response.Html([]byte("<H1>Default Output</H1>"), 200)
	}
}
