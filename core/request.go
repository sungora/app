package core

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// RequestUriParamsCompile
func RequestUriParamsCompile(postData map[string]interface{}) string {
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
	return q.Encode()
}

type requestHeader struct {
	authorizationBasic string
	contentType        string
	accept             string
}

// SetAuthorizationBasic установка заголовка Authorization
func (rh *requestHeader) SetAuthorizationBasic(login, passw string) {
	rh.authorizationBasic = "Basic " + base64.StdEncoding.EncodeToString([]byte(login+":"+passw))
}

// SetContentType установка заголовка ContentType
func (rh *requestHeader) SetContentType(contentType string) {
	rh.contentType = contentType
}

// func (rh *requestHeader) SetContentTypeJson() {
// 	rh.SetContentType("application/json")
// }

// SetAccept установка заголовка Accept
func (rh *requestHeader) SetAccept(accept string) {
	rh.accept = accept
}

// func (rh *requestHeader) SetAcceptJson() {
// 	rh.SetAccept("application/json")
// }

type request struct {
	url    string
	Header *requestHeader
}

// NewRequest создание запроса к внешнему ресурсу
func NewRequest(url string) *request {
	var r = new(request)
	r.url = url
	r.Header = &requestHeader{}
	return r
}

// GET запрос
func (r *request) GET(uri string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.request(http.MethodGet, uri, requestBody, responseBody)
}

// POST запрос
func (r *request) POST(uri string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.request(http.MethodPost, uri, requestBody, responseBody)
}

// PUT запрос
func (r *request) PUT(uri string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.request(http.MethodPut, uri, requestBody, responseBody)
}

// DELETE запрос
func (r *request) DELETE(uri string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.request(http.MethodDelete, uri, requestBody, responseBody)
}

// OPTIONS запрос
func (r *request) OPTIONS(uri string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.request(http.MethodOptions, uri, requestBody, responseBody)
}

func (r *request) request(method, uri string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	var url = r.url + uri
	var request *http.Request
	var data []byte
	body := new(bytes.Buffer)
	// Данные исходящего запроса
	if method == http.MethodPost || method == http.MethodPut {
		if data, err = json.Marshal(requestBody); err != nil {
			return
		}
		if _, err = body.Write(data); err != nil {
			return
		}
	} else if p, ok := requestBody.(map[string]interface{}); ok {
		url += "?" + RequestUriParamsCompile(p)
	}
	// Запрос
	if request, err = http.NewRequest(method, url, body); err == nil {
		// Заголовки
		if r.Header.authorizationBasic != "" {
			request.Header.Set("Authorization", r.Header.authorizationBasic)
		}
		if r.Header.contentType != "" {
			request.Header.Set("Content-Type", r.Header.contentType)
		}
		if r.Header.accept != "" {
			request.Header.Set("Accept", r.Header.accept)
		}
		//
		transCfg := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
		}
		c := http.Client{Transport: transCfg}
		if response, err = c.Do(request); err == nil {
			defer func() {
				err = response.Body.Close()
			}()
			var bodyResponse []byte
			bodyResponse, err = ioutil.ReadAll(response.Body)
			if err != nil {
				return
			}
			if r.Header.contentType == "application/json" {
				err = json.Unmarshal(bodyResponse, responseBody)
			}
			if response.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("%s:[%d]:%s", method, response.StatusCode, url))
			}
		}
	}
	return
}
