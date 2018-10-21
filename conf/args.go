// Инициализация параметров командной строки
package conf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// getCmdArgs Инициализация параметров командной строки
func GetCmdArgs() (mode string, err error) {
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	// - проверки
	if mode == `-h` || mode == `-help` || mode == `--help` {
		var str string
		str += "Usage of %s: %s [mode]\n"
		str += "\t mode: Режим запуска приложения\n"
		str += "\t\t install - Установка как сервиса в ОС\n"
		str += "\t\t uninstall - Удаление сервиса из ОС\n"
		str += "\t\t restart - Перезапуск ранее установленного сервиса\n"
		str += "\t\t start - Запуск ранее установленного сервиса\n"
		str += "\t\t stop - Остановка ранее установленного сервиса\n"
		str += "\t\t run - Прямой запуск (выход по 'Ctrl+C')\n"
		str += "\t\t если параметр опущен работает в режиме run\n"
		fmt.Fprintf(os.Stderr, str, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		return "", errors.New("Help startup request")
	}
	return
}
