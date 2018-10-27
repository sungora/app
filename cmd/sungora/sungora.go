// Служебные утилиты по работе со сборкой
// TODO Доработать трайсы функций (методы стираются)
package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Println("run - запуск приложения \n")
		return

	}

	switch os.Args[1] {
	case "run":
		var run = NewRun()
		run.Control()
	}
}
