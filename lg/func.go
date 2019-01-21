package lg

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"runtime"
	"time"

	"github.com/sungora/app/lg/message"
	"github.com/sungora/app/tool"
)

func Fatal(key interface{}, parametrs ...interface{}) {
	if config.Fatal == true {
		sendLog("fatal", "system", key, parametrs...)
	}
}
func FatalLogin(login string, key interface{}, parametrs ...interface{}) {
	if config.Fatal == true {
		sendLog("fatal", login, key, parametrs...)
	}
}
func Critical(key interface{}, parametrs ...interface{}) {
	if config.Critical == true {
		sendLog("critical", "system", key, parametrs...)
	}
}
func CriticalLogin(login string, key interface{}, parametrs ...interface{}) {
	if config.Critical == true {
		sendLog("critical", login, key, parametrs...)
	}
}
func Error(key interface{}, parametrs ...interface{}) {
	if config.Error == true {
		sendLog("error", "system", key, parametrs...)
	}
}
func ErrorLogin(login string, key interface{}, parametrs ...interface{}) {
	if config.Error == true {
		sendLog("error", login, key, parametrs...)
	}
}
func Warning(key interface{}, parametrs ...interface{}) {
	if config.Warning == true {
		sendLog("warning", "system", key, parametrs...)
	}
}
func WarningLogin(login string, key interface{}, parametrs ...interface{}) {
	if config.Warning == true {
		sendLog("warning", login, key, parametrs...)
	}
}
func Notice(key interface{}, parametrs ...interface{}) {
	if config.Notice == true {
		sendLog("notice", "system", key, parametrs...)
	}
}
func NoticeLogin(login string, key interface{}, parametrs ...interface{}) {
	if config.Notice == true {
		sendLog("notice", login, key, parametrs...)
	}
}
func Info(key interface{}, parametrs ...interface{}) {
	if config.Info == true {
		sendLog("info", "system", key, parametrs...)
	}
}
func InfoLogin(login string, key interface{}, parametrs ...interface{}) {
	if config.Info == true {
		sendLog("info", login, key, parametrs...)
	}
}
func Debug(key interface{}, parametrs ...interface{}) {
	if config.Debug == true {
		sendLog("debug", "system", key, parametrs...)
	}
}

func sendLog(level string, login string, key interface{}, parametrs ...interface{}) {
	msg := msg{}
	msg.Datetime = time.Now().Format(time.RFC850)
	msg.Level = level
	pc, _, line, ok := runtime.Caller(2)
	if ok == true {
		msg.LineNumber = line
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			msg.Action = fn.Name()
		}
	}
	msg.Service = tool.ServiceName
	msg.Login = login
	switch k := key.(type) {
	case error:
		msg.Message = k.Error()
	case int:
		msg.Message = message.GetMessage(k, parametrs...)
	case string:
		msg.Message = fmt.Sprintf(k, parametrs...)
	}
	if config.Traces == true {
		msg.Traces, _ = getTrace()
	}
	logCh <- msg
	// return msg.Message
}

// Dump all variables to STDOUT
func Dumper(idl ...interface{}) string {
	ret := dump(idl...)
	fmt.Print(ret.String())
	return ret.String()
}

// Dump all variables to bytes.Buffer
func dump(idl ...interface{}) bytes.Buffer {
	var buf bytes.Buffer
	var wr io.Writer

	wr = io.MultiWriter(&buf)
	for _, field := range idl {
		fset := token.NewFileSet()
		ast.Fprint(wr, fset, field, ast.NotNilFilter)
	}
	return buf
}
