package lg

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gopkg.in/sungora/app.v1/conf"
)

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
var config conf.Log
var timeLocation *time.Location
var serviceName string

func Start(c conf.Log, serviceName, timeZone string) (err error) {
	config = c
	config.ServiceName = serviceName
	// Инициализация временной зоны
	if loc, err := time.LoadLocation(timeZone); err == nil {
		timeLocation = loc
	} else {
		timeLocation = time.UTC
	}
	// Инициализация логирования в ФС
	sep := string(os.PathSeparator)
	if config.OutFilePath == "" {
		if dir, err := os.Getwd(); err == nil {
			config.OutFilePath = dir + sep + "logs" + sep + config.ServiceName + ".log"
		} else {
			return err
		}
	}
	if config.OutFile {
		os.MkdirAll(filepath.Dir(config.OutFilePath), 0777)
		fp, err = os.OpenFile(config.OutFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
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
			if config.OutHttp == true {
				saveHttp(msg)
			}
		}
		if fp != nil {
			fp.Close()
		}
		logChClose <- true
	}()
	return nil
}

func Wait() {
	close(logCh)
	<-logChClose
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
		tmp := strings.Split(infoList[i], "(")
		funcName := tmp[0]
		tmp = strings.Split(infoList[i+1], " ")
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
