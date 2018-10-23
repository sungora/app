package core

import (
	"net/http"
)

type ControllerJS struct {
	Request  *request
	Response *response
	Data     interface{}
}

func (self *ControllerJS) Init(w http.ResponseWriter, r *http.Request) {
	self.Request = newrequest(r)
	self.Response = newresponse(w)
}
func (self *ControllerJS) GET() (err error) {
	return
}
func (self *ControllerJS) POST() (err error) {
	return
}
func (self *ControllerJS) PUT() (err error) {
	return
}
func (self *ControllerJS) DELETE() (err error) {
	return
}
func (self *ControllerJS) OPTIONS() (err error) {
	return
}
func (self *ControllerJS) Render() {
	if self.Response.responseStatus {
		return
	}
	self.Response.Json200(self.Data)
}
