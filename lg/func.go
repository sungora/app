package lg

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"runtime"
	"time"

	"gopkg.in/sungora/app.v1/lg/message"
	"gopkg.in/sungora/app.v1/tool"
)

func Fatal(key interface{}, parametrs ...interface{}) error {
	if config.Fatal == true {
		return errors.New(sendLog("fatal", "system", key, parametrs...))
	}
	return errors.New("Fatal")
}
func FatalLogin(login string, key interface{}, parametrs ...interface{}) error {
	if config.Fatal == true {
		return errors.New(sendLog("fatal", login, key, parametrs...))
	}
	return errors.New("FatalLogin")
}
func Critical(key interface{}, parametrs ...interface{}) error {
	if config.Critical == true {
		return errors.New(sendLog("critical", "system", key, parametrs...))
	}
	return errors.New("Critical")
}
func CriticalLogin(login string, key interface{}, parametrs ...interface{}) error {
	if config.Critical == true {
		return errors.New(sendLog("critical", login, key, parametrs...))
	}
	return errors.New("CriticalLogin")
}
func Error(key interface{}, parametrs ...interface{}) error {
	if config.Error == true {
		return errors.New(sendLog("error", "system", key, parametrs...))
	}
	return errors.New("Error")
}
func ErrorLogin(login string, key interface{}, parametrs ...interface{}) error {
	if config.Error == true {
		return errors.New(sendLog("error", login, key, parametrs...))
	}
	return errors.New("ErrorLogin")
}
func Warning(key interface{}, parametrs ...interface{}) error {
	if config.Warning == true {
		return errors.New(sendLog("warning", "system", key, parametrs...))
	}
	return errors.New("Warning")
}
func WarningLogin(login string, key interface{}, parametrs ...interface{}) error {
	if config.Warning == true {
		return errors.New(sendLog("warning", login, key, parametrs...))
	}
	return errors.New("WarningLogin")
}
func Notice(key interface{}, parametrs ...interface{}) error {
	if config.Notice == true {
		return errors.New(sendLog("notice", "system", key, parametrs...))
	}
	return errors.New("Notice")
}
func NoticeLogin(login string, key interface{}, parametrs ...interface{}) error {
	if config.Notice == true {
		return errors.New(sendLog("notice", login, key, parametrs...))
	}
	return errors.New("NoticeLogin")
}
func Info(key interface{}, parametrs ...interface{}) error {
	if config.Info == true {
		return errors.New(sendLog("info", "system", key, parametrs...))
	}
	return errors.New("Info")
}
func InfoLogin(login string, key interface{}, parametrs ...interface{}) error {
	if config.Info == true {
		return errors.New(sendLog("info", login, key, parametrs...))
	}
	return errors.New("InfoLogin")
}
func Debug(key interface{}, parametrs ...interface{}) error {
	if config.Debug == true {
		return errors.New(sendLog("debug", "system", key, parametrs...))
	}
	return errors.New("Debug")
}

func sendLog(level string, login string, key interface{}, parametrs ...interface{}) string {
	msg := msg{}
	msg.Datetime = time.Now().In(tool.TimeLocation).Format("2006-01-02 15:04:05")
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
	case int:
		msg.Message = message.GetMessage(k, parametrs...)
	case string:
		msg.Message = fmt.Sprintf(k, parametrs...)
	}
	if config.Traces == true {
		msg.Traces, _ = getTrace()
	}
	logCh <- msg
	return msg.Message
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
