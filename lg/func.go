package lg

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Fatal error logging
func Fatal(key interface{}, parametrs ...interface{}) {
	if config.Fatal == true {
		sendLog("fatal", "system", key, parametrs...)
	}
}

// FatalLogin error logging
func FatalLogin(login string, key interface{}, parametrs ...interface{}) {
	if config.Fatal == true {
		sendLog("fatal", login, key, parametrs...)
	}
}

// Critical error logging
func Critical(key interface{}, parametrs ...interface{}) {
	if config.Critical == true {
		sendLog("critical", "system", key, parametrs...)
	}
}

// CriticalLogin error logging
func CriticalLogin(login string, key interface{}, parametrs ...interface{}) {
	if config.Critical == true {
		sendLog("critical", login, key, parametrs...)
	}
}

// Error error logging
func Error(key interface{}, parametrs ...interface{}) {
	if config.Error == true {
		sendLog("error", "system", key, parametrs...)
	}
}

// ErrorLogin error logging
func ErrorLogin(login string, key interface{}, parametrs ...interface{}) {
	if config.Error == true {
		sendLog("error", login, key, parametrs...)
	}
}

// Warning error logging
func Warning(key interface{}, parametrs ...interface{}) {
	if config.Warning == true {
		sendLog("warning", "system", key, parametrs...)
	}
}

// WarningLogin error logging
func WarningLogin(login string, key interface{}, parametrs ...interface{}) {
	if config.Warning == true {
		sendLog("warning", login, key, parametrs...)
	}
}

// Notice error logging
func Notice(key interface{}, parametrs ...interface{}) {
	if config.Notice == true {
		sendLog("notice", "system", key, parametrs...)
	}
}

// NoticeLogin error logging
func NoticeLogin(login string, key interface{}, parametrs ...interface{}) {
	if config.Notice == true {
		sendLog("notice", login, key, parametrs...)
	}
}

// Info error logging
func Info(key interface{}, parametrs ...interface{}) {
	if config.Info == true {
		sendLog("info", "system", key, parametrs...)
	}
}

// InfoLogin error logging
func InfoLogin(login string, key interface{}, parametrs ...interface{}) {
	if config.Info == true {
		sendLog("info", login, key, parametrs...)
	}
}

func sendLog(level string, login string, key interface{}, parametrs ...interface{}) {
	msg := msg{}
	msg.Datetime = time.Now().Format(time.RFC3339)
	msg.Level = level
	pc, _, line, ok := runtime.Caller(2)
	if ok == true {
		msg.LineNumber = line
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			msg.Action = fn.Name()
		}
	}
	msg.Service = config.ServiceName
	msg.Login = login
	switch k := key.(type) {
	case error:
		msg.Message = k.Error()
	case int:
		msg.Message = GetMessage(k, parametrs...)
	case string:
		msg.Message = fmt.Sprintf(k, parametrs...)
	}
	if config.Traces == true {
		msg.Traces, _ = getTrace()
	}
	component.logCh <- msg
	// return msg.Message
}

// Получение информаци о вызвавшем лог
func getTrace() (traces []trace, err error) {
	buf := make([]byte, 1<<16)
	i := runtime.Stack(buf, true)
	info := string(buf[:i])

	infoList := strings.Split(info, "\n")
	infoList = infoList[7:]

	for i := 0; i < len(infoList)-1; i += 2 {
		if infoList[i] == "" {
			break
		}
		if ok, _ := regexp.Match("/[gG]o/src/", []byte(infoList[i+1])); ok == true {
			break
		}
		funcName := infoList[i]
		tmp := strings.Split(infoList[i+1], " ")
		tmp = strings.Split(tmp[0], "go:")
		line, _ := strconv.Atoi(tmp[1])
		t := trace{
			FuncName:   funcName,
			FileName:   strings.TrimSpace(tmp[0]) + "go",
			LineNumber: line,
		}
		traces = append(traces, t)
	}
	return
}
