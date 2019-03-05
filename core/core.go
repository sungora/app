package core

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"os/signal"
	"reflect"
)

// Интерфейс компонентов приложения
type Componenter interface {
	// Инициализация компонентов приложения
	Init(cfg *ConfigRoot) (err error)
	// Запуск в работу компонентов приложения
	Start() (err error)
	// Завершение работы компонентов приложения
	Stop() (err error)
}

// ComponentReg Регистрация компонента приложения
func ComponentReg(com Componenter) {
	componentList = append(componentList, com)
}

var (
	componentList    []Componenter             // Срез зарегитрированных компонентов приложения
	chanelAppControl = make(chan os.Signal, 1) // Канал управления запуском и остановкой приложения
)

// Init Инициализация компонентов приложения
func Init() (code int) {
	var err error
	if len(componentList) == 0 {
		fmt.Fprintln(os.Stderr, "Ни одного компонента не зарегистрировано")
		return 1
	}
	componentList = append([]Componenter{component}, componentList...)

	// инициализация компонентов
	for i := 0; i < len(componentList); i++ {
		fmt.Fprintf(os.Stdout, "Init component %s\n", packageName(componentList[i]))
		if err = componentList[i].Init(Config); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
	}
	return
}

// Start Launch an application
func Start() (code int) {
	defer func() {
		chanelAppControl <- os.Interrupt
	}()
	var err error

	// 	завершение работы компонентов
	defer func() {
		for i := 0; i < len(componentList); i++ {
			fmt.Fprintf(os.Stdout, "Stop component %s\n", packageName(componentList[i]))
			if err = componentList[i].Stop(); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				code = 1
			}
		}
	}()

	// начало работы компонентов
	for i := 0; i < len(componentList); i++ {
		fmt.Fprintf(os.Stdout, "Start component %s\n", packageName(componentList[i]))
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

// Dump all variables to STDOUT
func Dumper(idl ...interface{}) string {
	ret := dump(idl...)
	fmt.Print(ret.String())
	return ret.String()
}

// dump all variables to bytes.Buffer
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
