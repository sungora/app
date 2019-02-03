package tool

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

type requestHeader struct {
	AuthorizationBasic string
	ContentType        string
	Accept             string
}

func (rh *requestHeader) SetAuthorizationBasic(login, passw string) {
	rh.AuthorizationBasic = "Basic " + base64.StdEncoding.EncodeToString([]byte(login+":"+passw))
}

func (rh *requestHeader) SetContentType(contentType string) {
	rh.ContentType = contentType
}

func (rh *requestHeader) SetAccept(accept string) {
	rh.Accept = accept
}

type request struct {
	url    string
	Header *requestHeader
}

func NewRequestJson(url string) *request {
	var r = new(request)
	r.url = url
	r.Header = &requestHeader{}
	r.Header.ContentType = "application/json"
	r.Header.Accept = "application/json"
	return r
}

func (r *request) GET(uri string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.request(http.MethodGet, uri, requestBody, responseBody)
}

func (r *request) POST(uri string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.request(http.MethodPost, uri, requestBody, responseBody)
}

func (r *request) PUT(uri string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.request(http.MethodPut, uri, requestBody, responseBody)
}

func (r *request) DELETE(uri string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.request(http.MethodDelete, uri, requestBody, responseBody)
}

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
		url += "?" + RequestGetParamsCompile(p)
	}
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}

	// Запрос
	if request, err = http.NewRequest(method, url, body); err == nil {
		// Заголовки
		if r.Header.AuthorizationBasic != "" {
			request.Header.Set("Authorization", r.Header.AuthorizationBasic)
		}
		request.Header.Set("Content-Type", r.Header.ContentType)
		c := http.Client{Transport: transCfg}
		if response, err = c.Do(request); err == nil {
			defer response.Body.Close()
			var bodyResponse []byte
			bodyResponse, err = ioutil.ReadAll(response.Body)
			if err != nil {
				return nil, err
			}
			if r.Header.ContentType == "application/json" {
				err = json.Unmarshal(bodyResponse, responseBody)
			}
			if response.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("%s:[%d]:%s", method, response.StatusCode, url))
			}
		}
	}
	return
}
