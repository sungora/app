package app

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

// Интерфейс компонентов приложения
type Componenter interface {
	// Запуск в работу компонентов приложения
	Start() (err error)
	// Завершение работы компонентов приложения
	Stop() (err error)
}

// ComponentAdd добавление компонента приложения
func ComponentAdd(com Componenter) {
	componentList = append(componentList, com)
}

var (
	componentList    []Componenter             // Срез зарегитрированных компонентов приложения
	chanelAppControl = make(chan os.Signal, 1) // Канал управления запуском и остановкой приложения
)

// Start Launch an application
func Start(IsStart *int8) (code int) {
	defer func() {
		chanelAppControl <- os.Interrupt
	}()
	var err error

	// начало работы компонентов
	for i := 0; i < len(componentList); i++ {
		fmt.Fprintf(os.Stdout, "Start component %s\n", packageName(componentList[i]))
		if err = componentList[i].Start(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			*IsStart = 1
			return 1
		}
	}

	// 	завершение работы компонентов
	defer func() {
		for i := len(componentList) - 1; -1 < i; i-- {
			fmt.Fprintf(os.Stdout, "Stop component %s\n", packageName(componentList[i]))
			if err = componentList[i].Stop(); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				*IsStart = 1
				code = 1
			}
		}
	}()
	*IsStart = 0

	// The correctness of the application is closed by a signal
	signal.Notify(chanelAppControl,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	<-chanelAppControl
	return
}

// Stop an application
func Stop() {
	chanelAppControl <- os.Interrupt
	<-chanelAppControl
}

// packageName Получение уникального имени пакета
func packageName(obj interface{}) string {
	var rt = reflect.TypeOf(obj)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return rt.PkgPath()
}
