package lg

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func saveHttp(m msg) {
	if bt, err := json.Marshal(m); err == nil {
		body := new(bytes.Buffer)
		if _, err := body.Write(bt); err == nil {
			if req, err := http.NewRequest("POST", config.OutHttpUrl, body); err == nil {
				c := http.Client{}
				if resp, err := c.Do(req); err == nil {
					resp.Body.Close()
				}
			}
		}
	}
}
