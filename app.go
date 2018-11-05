// Служебные утилиты по работе со сборкой
// TODO Доработать трайсы функций (методы стираются)
package main

import (
	"fmt"
	"os"
)

const VERSION = "v1"

func main() {

	if len(os.Args) == 1 {
		fmt.Println("run - сборка и запуск приложения")
		fmt.Println("new nameApp - новое приложение")
		return
	}

	switch os.Args[1] {
	case "run":
		var run = NewRun()
		run.Control()
	case "new":
		if len(os.Args) != 3 {
			fmt.Println("имя приложения не задано")
			return
		}
		var app = newApp(os.Args[2], VERSION)
		if app == nil {
			fmt.Println("бланк приложения не найден")
			return
		}
		if err := app.Copy(); err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("OK")
		}
	}
}
