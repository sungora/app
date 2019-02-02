package lg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/sungora/app/lg/message"
	"github.com/sungora/app/startup"
)

// init регистрация компонента в приложении
func init() {
	component = new(componentTyp)
	startup.SetComponent(component)
}

// компонент
type componentTyp struct {
	fp         *os.File  // запись логов в файл
	logCh      chan msg  // канал чтения и обработки логов
	logChClose chan bool // канал управления закрытием работы
}

var (
	config    *configMain   // конфигурация
	component *componentTyp // компонент
)

// Init инициализация компонента в приложении
func (comp *componentTyp) Init() (err error) {
	sep := string(os.PathSeparator)
	config = new(configMain)

	// техническое имя приложения
	if ext := filepath.Ext(os.Args[0]); ext != "" {
		sl := strings.Split(filepath.Base(os.Args[0]), ext)
		config.ServiceName = sl[0]
	} else {
		config.ServiceName = filepath.Base(os.Args[0])
	}

	// диреткория логов приложения
	dirWork, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	config.dirLog = dirWork + sep + "log"
	var fi os.FileInfo
	if fi, err = os.Stat(config.dirLog); err != nil {
		if err = os.MkdirAll(config.dirLog, 0700); err != nil {
			return
		}
	} else if fi.IsDir() == false {
		return errors.New("не правильная директория логов\n" + config.dirLog)
	}

	// читаем конфигурацию
	path := dirWork + sep + "config" + sep + config.ServiceName + ".toml"
	if _, err = toml.DecodeFile(path, &config); err != nil {
		return
	}

	// читаем шаблоны сообщений логов
	msgTmp := make(map[string]string)
	path = dirWork + sep + "config" + sep + config.ServiceName + "_lg.toml"
	if _, err := toml.DecodeFile(path, &msgTmp); err != nil {
		fmt.Fprintln(os.Stdout, err.Error())
	} else {
		for codeStr, msg := range msgTmp {
			if code, err := strconv.Atoi(codeStr); err == nil {
				message.SetMessage(code, msg)
			}
		}
	}

	comp.logCh = make(chan msg, 10000) // канал чтения и обработки логов
	comp.logChClose = make(chan bool)  // канал управления закрытием работы

	return
}

// Start запуск компонента в работу
func (comp *componentTyp) Start() (err error) {
	if config.Lg.OutFile {
		if comp.fp, err = os.OpenFile(config.dirLog+"/"+config.ServiceName+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600); err != nil {
			return
		}
	}
	go func() {
		for msg := range comp.logCh {
			if config.Lg.OutStd == true {
				saveStdout(msg)
			}
			if config.Lg.OutFile == true {
				saveFile(msg)
			}
			if config.Lg.OutHttp != "" {
				saveHttp(msg)
			}
		}
		comp.logChClose <- true
	}()
	return
}

// Stop завершение работы компонента
func (comp *componentTyp) Stop() (err error) {
	if comp.fp != nil {
		err = comp.fp.Close()
	}
	close(comp.logCh)
	<-comp.logChClose
	return
}
