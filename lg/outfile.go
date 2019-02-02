package lg

import "encoding/json"

func saveFile(m msg) {
	var logLine string
	if bt, err := json.Marshal(m); err == nil {
		logLine = string(bt)
	} else {
		return
	}
	if component.fp != nil {
		component.fp.WriteString(logLine + "\n")
	}
}
