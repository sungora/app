package lg

import (
	"encoding/json"
	"os"
)

var fp *os.File

func saveFile(m msg) {
	var logLine string
	if bt, err := json.Marshal(m); err == nil {
		logLine = string(bt)
	} else {
		return
	}
	if fp != nil {
		fp.WriteString(logLine + "\n")
	}
}
