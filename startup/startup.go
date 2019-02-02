package startup // import "github.com/sungora/app/startup"

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
)

type ComponentFace interface {
	// Инициализация компонентов приложения
	Init() (err error)
	// Запуск в работу компонентов приложения
	Start() (err error)
	// Завершение работы компонентов приложения
	Stop() (err error)
}

var componentList []ComponentFace

func SetComponent(com ComponentFace) {
	componentList = append(componentList, com)
}

var (
	chanelAppControl = make(chan os.Signal, 1) // Канал управления запуском и остановкой приложения
)

// Start Launch an application
func Start() (code int) {

	defer func() {
		chanelAppControl <- os.Interrupt
	}()
	var err error

	if len(componentList) == 0 {
		fmt.Fprintln(os.Stderr, "Ни одного компонента не зарегистрировано")
		return 1
	}

	// инициализация
	for i := 0; i < len(componentList); i++ {
		fmt.Fprintf(os.Stdout, "Init component %s\n", packageName(componentList[i]))
		if err = componentList[i].Init(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}

	}

	// 	завершение работы
	defer func() {
		for i := 0; i < len(componentList); i++ {
			if err = componentList[i].Stop(); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				code = 1
			}
		}
	}()

	// начало в работы
	for i := 0; i < len(componentList); i++ {
		if err = componentList[i].Start(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
	}

	// The correctness of the application is closed by a signal
	signal.Notify(chanelAppControl, os.Interrupt)
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
