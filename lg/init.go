package lg

import (
	"errors"
	"os"
	"path/filepath"
)

// init регистрация компонента в приложении
// func init() {
// 	component = new(Component)
// 	core.ComponentReg(component)
// }

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
func Init(cfg *Config) (com *Component, err error) {

	config = cfg
	component = new(Component)

	// sep := string(os.PathSeparator)
	// config = new(Config)
	// config.ServiceName = cfg.ServiceName

	// диреткория логов приложения
	dir := filepath.Dir(os.Args[0])

	var fi os.FileInfo
	if fi, err = os.Stat(dir); err != nil {
		if err = os.MkdirAll(dir, 0700); err != nil {
			return
		}
	} else if fi.IsDir() == false {
		err = errors.New("не правильная директория логов\n" + dir)
	}

	// читаем конфигурацию
	// path := cfg.DirConfig + sep + cfg.ServiceName + ".toml"
	// if _, err = toml.DecodeFile(path, config); err != nil {
	// 	return
	// }

	// читаем шаблоны сообщений логов
	// msgTmp := make(map[string]string)
	// path = cfg.DirConfig + sep + cfg.ServiceName + "_lg.toml"
	// if _, err := toml.DecodeFile(path, &msgTmp); err != nil {
	// 	fmt.Fprintln(os.Stdout, err.Error())
	// } else {
	// 	for codeStr, msg := range msgTmp {
	// 		if code, err := strconv.Atoi(codeStr); err == nil {
	// 			message.SetMessage(code, msg)
	// 		}
	// 	}
	// }

	return component, nil
}

// Start запуск компонента в работу
func (comp *Component) Start() (err error) {

	comp.logCh = make(chan msg, 10000) // канал чтения и обработки логов
	comp.logChClose = make(chan bool)  // канал управления закрытием работы

	if config.Lg.OutFile != "" {
		if comp.fp, err = os.OpenFile(config.Lg.OutFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600); err != nil {
			return
		}
	}
	go func() {
		for msg := range comp.logCh {
			if config.Lg.OutStd == true {
				saveStdout(msg)
			}
			if config.Lg.OutFile != "" {
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
func (comp *Component) Stop() (err error) {
	if comp.fp != nil {
		err = comp.fp.Close()
	}
	close(comp.logCh)
	<-comp.logChClose
	return
}
