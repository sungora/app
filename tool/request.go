package tool

import (
	"bytes"
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
	contentType        string
	authorizationBasic string
}

func (rh *requestHeader) SetContentType(contentType string) {
	rh.contentType = contentType
}
func (rh *requestHeader) AuthorizationBasic(login, passw string) {
	rh.authorizationBasic = "Basic " + base64.StdEncoding.EncodeToString([]byte(login+":"+passw))
}

type request struct {
	url    string
	Header *requestHeader
}

func NewRequestJson(url string) *request {
	var r = new(request)
	r.url = url
	r.Header = &requestHeader{}
	r.Header.contentType = "application/json"
	return r
}

func (r *request) GET(requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.request("GET", requestBody, responseBody)
}

func (r *request) POST(requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.request("POST", requestBody, responseBody)
}

func (r *request) PUT(requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.request("PUT", requestBody, responseBody)
}

func (r *request) DELETE(requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.request("DELETE", requestBody, responseBody)
}

func (r *request) OPTIONS(requestBody, responseBody interface{}) (response *http.Response, err error) {
	return r.request("OPTIONS", requestBody, responseBody)
}

func (r *request) request(method string, requestBody, responseBody interface{}) (response *http.Response, err error) {
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
			var bodyResponse []byte
			bodyResponse, err = ioutil.ReadAll(response.Body)
			if err != nil {
				return nil, err
			}
			if r.Header.contentType == "application/json" {
				err = json.Unmarshal(bodyResponse, responseBody)
			}
			if response.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("%s:[%d]:%s", method, response.StatusCode, r.url))
			}
		}
	}
	return
}
