package lg

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	config    *Config    // конфигурация
	component *Component // компонент
)

// компонент
type Component struct {
	fp         *os.File  // запись логов в файл
	logCh      chan msg  // канал чтения и обработки логов
	logChClose chan bool // канал управления закрытием работы
}

// Init инициализация компонента в приложении
func Init(cfg *Config, serviceName string) (com *Component, err error) {
	config = cfg
	component = new(Component)
	// сервис
	if config.ServiceName == "" {
		config.ServiceName = serviceName
	}
	// диреткория логов приложения
	var dir string
	if config.OutFile == "" {
		sep := string(os.PathSeparator)
		dir = filepath.Dir(filepath.Dir(os.Args[0]))
		if dir == "." {
			dir, _ = os.Getwd()
			dir = filepath.Dir(dir)
		}
		dir += "/log"
		config.OutFile = dir + sep + serviceName + ".log"
	} else {
		dir = filepath.Dir(config.OutFile)
	}
	var fi os.FileInfo
	if fi, err = os.Stat(dir); err != nil {
		if err = os.MkdirAll(dir, 0700); err != nil {
			return
		}
	} else if fi.IsDir() == false {
		return nil, errors.New("не правильная директория логов\n" + dir)
	}
	//
	return component, nil
}

// Start запуск компонента в работу
func (comp *Component) Start() (err error) {

	comp.logCh = make(chan msg, 10000) // канал чтения и обработки логов
	comp.logChClose = make(chan bool)  // канал управления закрытием работы

	if config.OutFile != "" {
		if comp.fp, err = os.OpenFile(config.OutFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600); err != nil {
			return
		}
	}
	go func() {
		for msg := range comp.logCh {
			if config.OutStd == true {
				saveStdout(msg)
			}
			if config.OutFile != "" {
				saveFile(msg)
			}
			if config.OutHttp != "" {
				saveHttp(msg)
			}
		}
		comp.logChClose <- true
	}()
	return
}

// Stop завершение работы компонента
func (comp *Component) Stop() (err error) {
	if comp.fp != nil {
		err = comp.fp.Close()
	}
	close(comp.logCh)
	<-comp.logChClose
	return
}
