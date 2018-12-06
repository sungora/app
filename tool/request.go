package tool

import (
	"bytes"
	"encoding/json"
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

func NewRequestGET(url string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return NewRequest(url, "GET", requestBody, responseBody)
}

func NewRequestPOST(url string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return NewRequest(url, "POST", requestBody, responseBody)
}

func NewRequestPUT(url string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return NewRequest(url, "PUT", requestBody, responseBody)
}

func NewRequestDELETE(url string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return NewRequest(url, "DELETE", requestBody, responseBody)
}

func NewRequestOPTIONS(url string, requestBody, responseBody interface{}) (response *http.Response, err error) {
	return NewRequest(url, "OPTIONS", requestBody, responseBody)
}

func NewRequest(url, method string, requestBody, responseBody interface{}) (response *http.Response, err error) {
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
		c := http.Client{}
		if response, err = c.Do(request); err == nil {
			defer response.Body.Close()
			bodyResponse, _ := ioutil.ReadAll(response.Body)
			json.Unmarshal(bodyResponse, responseBody)
		}
	}
	return
}
