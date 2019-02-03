package lg

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func saveHttp(m msg) {
	if data, err := json.Marshal(m); err == nil {
		body := new(bytes.Buffer)
		if _, err := body.Write(data); err == nil {
			if request, err := http.NewRequest("POST", config.OutHttp, body); err == nil {
				request.Header.Set("Content-Type", "application/json")
				c := http.Client{}
				if resp, err := c.Do(request); err == nil {
					resp.Body.Close()
				}
			}
		}
	}
}
