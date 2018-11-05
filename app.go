// Служебные утилиты по работе со сборкой
// TODO Доработать трайсы функций (методы стираются)
package main

import (
	"fmt"
	"os"
)

const VERSION = "v1"

func main() {

	if len(os.Args) < 3 {
		fmt.Println("run nameApp - авто-сборка и авто-запуск приложения")
		fmt.Println("new nameApp - создать новое приложение")
		return
	}

	switch os.Args[1] {
	case "run":
		var run = NewRun(os.Args[2])
		run.Control()
	case "new":
		var app = newApp(os.Args[2], VERSION)
		if app == nil {
			fmt.Println("бланк приложения не найден")
			return
		}
		if err := app.New(); err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("OK")
		}
	case "blank":
		var app = newApp(os.Args[2], VERSION)
		if app == nil {
			fmt.Println("бланк приложения не найден")
			return
		}
		if err := app.Blank(); err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("OK")
		}
	}
}
