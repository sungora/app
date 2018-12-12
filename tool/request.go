package tool

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"errors"
)

func RequestGetParamsCompile(postData map[string]interface{}) string {
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

func NewRequestGET(url string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return newRequest(url, "GET", requestBody, responseBody)
}

func NewRequestPOST(url string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return newRequest(url, "POST", requestBody, responseBody)
}

func NewRequestPUT(url string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return newRequest(url, "PUT", requestBody, responseBody)
}

func NewRequestDELETE(url string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return newRequest(url, "DELETE", requestBody, responseBody)
}

func NewRequestOPTIONS(url string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return newRequest(url, "OPTIONS", requestBody, responseBody)
}

func newRequest(url, method string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	var request *http.Request
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
		url += "?" + RequestGetParamsCompile(p)
	}
	//
	if request, err = http.NewRequest(method, url, body); err == nil {
		request.Header.Set("Content-Type", "application/json")
		en := base64.StdEncoding.EncodeToString([]byte("Inventory:ByFIPhwipuZ7fthaAq3DnjBEJQiS6sG"))
		request.Header.Set("Authorization", "Basic "+en)
		// request.Header.Set("Authorization", "Basic Inventory:ByFIPhwipuZ7fthaAq3DnjBEJQiS6sG")
		c := http.Client{}
		if response, err = c.Do(request); err == nil {
			defer response.Body.Close()
			bodyResponse, _ := ioutil.ReadAll(response.Body)
			json.Unmarshal(bodyResponse, responseBody)
		}
	}
	return
}

// /////

type requestHeader struct {
	contentType        string
	authorizationBasic string
}

func (rh *requestHeader) ContentType(contentType string) {
	rh.contentType = contentType
}
func (rh *requestHeader) AuthorizationBasic(login, passw string) {
	rh.authorizationBasic = "Basic " + base64.StdEncoding.EncodeToString([]byte(login+":"+passw))
}

type request struct {
	url    string
	Header *requestHeader
}

func NewRequest(url string) *request {
	var r = new(request)
	r.url = url
	r.Header = &requestHeader{}
	return r
}

func (r *request) GET(requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.qwerty("GET", requestBody, responseBody)
}

func (r *request) POST(requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.qwerty("POST", requestBody, responseBody)
}

func (r *request) PUT(requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.qwerty("PUT", requestBody, responseBody)
}

func (r *request) DELETE(requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.qwerty("DELETE", requestBody, responseBody)
}

func (r *request) OPTIONS(requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.qwerty("OPTIONS", requestBody, responseBody)
}

func (r *request) qwerty(method string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	var request *http.Request
	body := new(bytes.Buffer)
	// Данные исходящего запроса
	if method == "POST" || method == "PUT" {
		var data []byte
		if data, err = json.Marshal(requestBody); err != nil {
			return
		}
		if _, err = body.Write(data); err != nil {
			return
		}
	} else if p, ok := requestBody.(map[string]interface{}); ok {
		r.url += "?" + RequestGetParamsCompile(p)
	}
	// Запрос
	if request, err = http.NewRequest(method, r.url, body); err == nil {
		// Заголовки
		if r.Header.authorizationBasic != "" {
			request.Header.Set("Authorization", r.Header.authorizationBasic)
		}
		request.Header.Set("Content-Type", r.Header.contentType)
		c := http.Client{}
		if response, err = c.Do(request); err == nil {
			defer response.Body.Close()
			//
			if r.Header.contentType == "application/json" {
				bodyResponse, _ := ioutil.ReadAll(response.Body)
				err = json.Unmarshal(bodyResponse, responseBody)
			}
			//
			if response.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("%d:[%s]:%s", method, response.StatusCode, r.url))
			}
		}
	}
	return
}
