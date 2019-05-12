package lg

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	config    *Config    // конфигурация
	component *Component // компонент
)

// Init инициализация компонента в приложении
func Init(cfg *Config) (com *Component, err error) {
	config = cfg
	component = new(Component)
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
		serviceName := filepath.Base(os.Args[0])
		if filepath.Ext(serviceName) != "" {
			serviceName = strings.Split(serviceName, filepath.Ext(serviceName))[0]
		}
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

func GetConfig() *Config {
	return config
}
