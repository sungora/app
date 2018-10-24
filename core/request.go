package core

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"gopkg.in/sungora/app.v1/lg"
)

type request struct {
	Request *http.Request
	Params  map[string][]string
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
	// l := strings.Split(r.Header.Get("Content-Type"), ";")
	// self.ContentType = l[0]
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

func RequestGetParams(postData map[string]interface{}) string {
	u := new(url.URL)
	q := u.Query()
	for k, v := range postData {
		switch v1 := v.(type) {
		case uint64:
			q.Add(k, strconv.FormatUint(v1, 10))
		case int64:
			q.Add(k, strconv.FormatInt(v1, 10))
		case int:
			q.Add(k, strconv.Itoa(v1))
		case float64:
			q.Add(k, strconv.FormatFloat(v1, 'f', -1, 64))
		case bool:
			q.Add(k, strconv.FormatBool(v1))
		case string:
			q.Add(k, v1)
		}
	}
	// query = strings.TrimLeft(query, "&")
	return q.Encode()
}

func NewRequestGET(url string, requestBody, responseBody interface{}) (err error) {
	return NewRequest(url, "GET", requestBody, responseBody)
}

func NewRequestPOST(url string, requestBody, responseBody interface{}) (err error) {
	return NewRequest(url, "POST", requestBody, responseBody)
}

func NewRequestPUT(url string, requestBody, responseBody interface{}) (err error) {
	return NewRequest(url, "PUT", requestBody, responseBody)
}

func NewRequestDELETE(url string, requestBody, responseBody interface{}) (err error) {
	return NewRequest(url, "DELETE", requestBody, responseBody)
}

func NewRequestOPTIONS(url string, requestBody, responseBody interface{}) (err error) {
	return NewRequest(url, "OPTIONS", requestBody, responseBody)
}

func NewRequest(url, method string, requestBody, responseBody interface{}) (err error) {
	var request *http.Request
	var response *http.Response
	body := new(bytes.Buffer)
	//
	if method == "POST" || method == "PUT" {
		var data []byte
		if data, err = json.Marshal(requestBody); err != nil {
			return
		}
		if _, err = body.Write(data); err != nil {
			return
		}
	} else if p, ok := requestBody.(map[string]interface{}); ok {
		url += "?" + RequestGetParams(p)
	}
	//
	if request, err = http.NewRequest(method, url, body); err == nil {
		request.Header.Set("Content-Type", "application/json")
		c := http.Client{}
		if response, err = c.Do(request); err == nil {
			defer response.Body.Close()
			bodyResponse, _ := ioutil.ReadAll(response.Body)
			json.Unmarshal(bodyResponse, responseBody)
		}
	}
	return err
}

func _NewRequest(url, method string, requestBody, responseBody interface{}) (err error) {
	var request *http.Request
	var response *http.Response
	data, err := json.Marshal(requestBody)
	body := new(bytes.Buffer)
	if _, err = body.Write(data); err == nil {
		if request, err = http.NewRequest(method, url, body); err == nil {
			request.Header.Set("Content-Type", "application/json")
			c := http.Client{}
			if response, err = c.Do(request); err == nil {
				defer response.Body.Close()
				bodyResponse, _ := ioutil.ReadAll(response.Body)
				json.Unmarshal(bodyResponse, responseBody)
			}
		}
	}
	return err
}
