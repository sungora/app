package app

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
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

// Start Launch an application (начало работы компонентов)
func Start() (err error) {
	for i := 0; i < len(componentList); i++ {
		fmt.Fprintf(os.Stdout, "Start component %s\n", packageName(componentList[i]))
		if err = componentList[i].Start(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
	}
	return
}

// Stop an application (завершение работы компонентов)
func Stop() {
	for i := len(componentList) - 1; -1 < i; i-- {
		fmt.Fprintf(os.Stdout, "Stop component %s\n", packageName(componentList[i]))
		if err := componentList[i].Stop(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
	}
	return
}

// StartLock Launch an application
func StartLock() (err error) {
	defer func() {
		chanelAppControl <- os.Interrupt
	}()

	if err = Start(); err != nil {
		return
	}
	defer Stop()

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

// StopLock an application
func StopLock() {
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
		_ = ast.Fprint(wr, fset, field, ast.NotNilFilter)
	}
	return buf
}
