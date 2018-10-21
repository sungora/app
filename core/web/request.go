package web

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/sungora/app.v1/lg"
)

type request struct {
	Request     *http.Request
	Params      map[string][]string
	ContentType string
}

func newrequest(r *http.Request) *request {
	// var uriParams, _ = url.ParseQuery(r.URL.Query().Encode())
	// r.Header.Get("Content-Type")
	self := new(request)
	self.Request = r
	self.Params, _ = url.ParseQuery(r.URL.Query().Encode())
	r.ParseForm()
	for i, v := range r.Form {
		self.Params[i] = v
	}
	//
	l := strings.Split(r.Header.Get("Content-Type"), ";")
	self.ContentType = l[0]
	// switch l[0] {
	// case "text/plain":
	// case "application/x-www-form-urlencoded":
	// case "multipart/form-data":
	// case "application/json":
	// }
	return self
}

func (self *request) BodyDecode(object interface{}) error {
	if body, err := ioutil.ReadAll(self.Request.Body); err == nil {
		if err = json.Unmarshal(body, object); err != nil {
			return lg.Error(err.Error())
		}
	} else {
		return lg.Error(err.Error())
	}
	return nil
}

func (self *request) Send(url, method string, objSet, objGet interface{}) (err error) {
	var request *http.Request
	var response *http.Response
	data, err := json.Marshal(objSet)
	body := new(bytes.Buffer)
	if _, err = body.Write(data); err == nil {
		if request, err = http.NewRequest(method, url, body); err == nil {
			request.Header.Set("Content-Type", "application/json")
			c := http.Client{}
			if response, err = c.Do(request); err == nil {
				defer response.Body.Close()
				bodyResponse, _ := ioutil.ReadAll(response.Body)
				json.Unmarshal(bodyResponse, objGet)
			}
		}
	}
	return err
}
