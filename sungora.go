// Служебные утилиты по работе со сборкой
// TODO Доработать трайсы функций (методы стираются)
package main

import (
	"fmt"
	"gopkg.in/sungora/app.v1/utils"
	"os"
)

func main() {
	os.Chdir("/home/konstantin/go/src/accounter")

	if len(os.Args) == 1 {
		fmt.Println("run - запуск приложения \n")
		return

	}

	switch os.Args[1] {
	case "run":
		var run = utils.NewRun()
		run.Control()
	}
}
