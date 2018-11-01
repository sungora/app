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

	"gopkg.in/sungora/app.v1/conf"
	"gopkg.in/sungora/app.v1/lg/message"
)

func Fatal(key interface{}, messages ...interface{}) error {
	if config.Fatal == true {
		return errors.New(sendLog("fatal", "system", key, messages...))
	}
	return errors.New("Fatal")
}
func FatalLogin(login string, key interface{}, messages ...interface{}) error {
	if config.Fatal == true {
		return errors.New(sendLog("fatal", login, key, messages...))
	}
	return errors.New("FatalLogin")
}
func Critical(key interface{}, messages ...interface{}) error {
	if config.Critical == true {
		return errors.New(sendLog("critical", "system", key, messages...))
	}
	return errors.New("Critical")
}
func CriticalLogin(login string, key interface{}, messages ...interface{}) error {
	if config.Critical == true {
		return errors.New(sendLog("critical", login, key, messages...))
	}
	return errors.New("CriticalLogin")
}
func Error(key interface{}, messages ...interface{}) error {
	if config.Error == true {
		return errors.New(sendLog("error", "system", key, messages...))
	}
	return errors.New("Error")
}
func ErrorLogin(login string, key interface{}, messages ...interface{}) error {
	if config.Error == true {
		return errors.New(sendLog("error", login, key, messages...))
	}
	return errors.New("ErrorLogin")
}
func Warning(key interface{}, messages ...interface{}) error {
	if config.Warning == true {
		return errors.New(sendLog("warning", "system", key, messages...))
	}
	return errors.New("Warning")
}
func WarningLogin(login string, key interface{}, messages ...interface{}) error {
	if config.Warning == true {
		return errors.New(sendLog("warning", login, key, messages...))
	}
	return errors.New("WarningLogin")
}
func Notice(key interface{}, messages ...interface{}) error {
	if config.Notice == true {
		return errors.New(sendLog("notice", "system", key, messages...))
	}
	return errors.New("Notice")
}
func NoticeLogin(login string, key interface{}, messages ...interface{}) error {
	if config.Notice == true {
		return errors.New(sendLog("notice", login, key, messages...))
	}
	return errors.New("NoticeLogin")
}
func Info(key interface{}, messages ...interface{}) error {
	if config.Info == true {
		return errors.New(sendLog("info", "system", key, messages...))
	}
	return errors.New("Info")
}
func InfoLogin(login string, key interface{}, messages ...interface{}) error {
	if config.Info == true {
		return errors.New(sendLog("info", login, key, messages...))
	}
	return errors.New("InfoLogin")
}
func Debug(key interface{}, messages ...interface{}) error {
	if config.Debug == true {
		return errors.New(sendLog("debug", "system", key, messages...))
	}
	return errors.New("Debug")
}

func sendLog(level string, login string, key interface{}, messages ...interface{}) string {
	msg := msg{}
	msg.Datetime = time.Now().In(conf.TimeLocation).Format("2006-01-02 15:04:05")
	msg.Level = level
	pc, _, line, ok := runtime.Caller(2)
	if ok == true {
		msg.LineNumber = line
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			msg.Action = fn.Name()
		}
	}
	msg.Service = serviceName
	msg.Login = login
	switch k := key.(type) {
	case int:
		msg.Message = message.GetMessage(k, messages...)
	case string:
		msg.Message = fmt.Sprintf(k, messages...)
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
