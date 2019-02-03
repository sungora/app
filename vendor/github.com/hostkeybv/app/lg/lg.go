package lg

import (
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/hostkeybv/app/tool"
)

type Config struct {
	Info     bool
	Notice   bool
	Warning  bool
	Error    bool
	Critical bool
	Fatal    bool
	Debug    bool
	Traces   bool
	OutStd   bool
	OutFile  bool
	OutHttp  string // url куда отправляются логи
}

type msg struct {
	Datetime   string
	Level      string
	LineNumber int
	Action     string
	Service    string
	Login      string
	Message    string
	Traces     []trace
}

var logCh = make(chan msg, 10000)
var logChClose = make(chan bool)
var config Config

func Init(c Config) (err error) {
	config = c

	// Инициализация логирования в ФС
	if tool.DirLog != "" {
		os.MkdirAll(tool.DirLog, 0777)
		fp, err = os.OpenFile(tool.DirLog+"/"+tool.ServiceName+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return
		}
	}

	// Create a PID file and lock on record, control run one copy of the application
	os.MkdirAll(tool.DirPid, 0777)
	if err = pidFileCreate(tool.DirPid + "/pid"); err != nil {
		return
	}

	//
	go func() {
		for msg := range logCh {
			if config.OutStd == true {
				saveStdout(msg)
			}
			if config.OutFile == true {
				saveFile(msg)
			}
			if config.OutHttp != "" {
				saveHttp(msg)
			}
		}
		if fp != nil {
			fp.Close()
		}
		logChClose <- true
	}()
	return
}

func Wait() {
	close(logCh)
	<-logChClose
	pidFileUnlock()
}

type trace struct {
	FuncName   string // Название функции
	FileName   string // Имя исходного файла
	LineNumber int    // Номер строки внутри функции
}

// Получение информаци о вызвавшем лог
func getCallInfo(level int) (FuncName string, LineNumber int, FileName string) {
	pc, file, line, ok := runtime.Caller(level)
	if ok == true {
		LineNumber = line
		FileName = file
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			FuncName = fn.Name()
		}
	}
	return
}

// Получение информаци о вызвавшем лог
func getTrace() (traces []trace, err error) {
	buf := make([]byte, 1<<16)
	i := runtime.Stack(buf, true)
	info := string(buf[:i])

	infoList := strings.Split(info, "\n")
	infoList = infoList[7:]

	// fmt.Println(infoList[0])
	// fmt.Println(info)

	for i := 0; i < len(infoList)-1; i += 2 {
		if infoList[i] == "" {
			break
		}
		if ok, _ := regexp.Match("/[gG]o/src/", []byte(infoList[i+1])); ok == true {
			break
		}
		// tmp := strings.Split(infoList[i], "(")
		// funcName := tmp[0]
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
