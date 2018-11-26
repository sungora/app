// Управление запуском и остановкой приложения
package service

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"gopkg.in/sungora/app.v1/core"
	"gopkg.in/sungora/app.v1/lg"
	"gopkg.in/sungora/app.v1/tool"
	"gopkg.in/sungora/app.v1/workflow"
)

type config struct {
	Core       core.ConfigTyp
	Mysql      configMysql
	Postgresql configPostgresql
	Log        lg.Config
	Workflow   workflow.Config
}

type configMysql struct {
	Host     string // протокол, хост и порт подключения
	Name     string // Имя базы данных
	Login    string // Логин к базе данных
	Password string // Пароль к базе данных
	Charset  string // Кодировка данных (utf-8 - по умолчанию)
}

type configPostgresql struct {
	Host     string // Хост базы данных (localhost - по умолчанию)
	Port     int64  // Порт подключения по протоколу tcp/ip (3306 по умолчанию)
	Name     string // Имя базы данных
	Login    string // Логин к базе данных
	Password string // Пароль к базе данных
	Charset  string // Кодировка данных (utf-8 - по умолчанию)
}

// Каналы управления запуском и остановкой приложения
var (
	chanelAppStop    = make(chan os.Signal, 1)
	chanelAppControl = make(chan os.Signal, 1)
)

// Start Launch an application
func Start() (code int) {
	defer func() { // контроль завершение работы приложения
		chanelAppStop <- os.Interrupt
	}()
	var (
		err   error
		store net.Listener
	)

	// config
	var configApp *config
	path := tool.DirConfig + string(os.PathSeparator) + "main.toml"
	if _, err = toml.DecodeFile(path, &configApp); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	core.Config = &configApp.Core
	core.Config.SessionTimeout *= time.Second

	// Инициализация временной зоны
	if loc, err := time.LoadLocation(configApp.Core.TimeZone); err == nil {
		tool.TimeLocation = loc
	} else {
		tool.TimeLocation = time.UTC
	}

	// logs
	if err = lg.Start(configApp.Log); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	defer lg.Wait()

	// core
	switch configApp.Core.DriverDB {
	case "mysql":
		if core.DB, err = gorm.Open("mysql", fmt.Sprintf(
			"%s:%s@%s/%s?charset=%s&parseTime=True&loc=Local&timeout=3s",
			configApp.Mysql.Login,
			configApp.Mysql.Password,
			configApp.Mysql.Host,
			configApp.Mysql.Name,
			configApp.Mysql.Charset,
		)); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
		defer core.DB.Close()
	case "postgresql":
		if core.DB, err = gorm.Open("postgres", fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s",
			configApp.Postgresql.Host,
			configApp.Postgresql.Port,
			configApp.Postgresql.Login,
			configApp.Postgresql.Name,
			configApp.Postgresql.Password,
		)); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
		defer core.DB.Close()
	}

	// session
	if 0 < configApp.Core.SessionTimeout {
		core.SessionGC()
	}

	// workflow
	if configApp.Workflow.IsWorkflow == true {
		if err = workflow.Start(configApp.Workflow); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
		defer workflow.Wait()
	}

	// service - application
	if store, err = newHttp(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	defer store.Close()
	fmt.Fprintf(os.Stdout, "service start success: http://%s:%d\n", core.Config.Host, core.Config.Port)

	// The correctness of the application is closed by a signal
	signal.Notify(chanelAppControl, os.Interrupt)
	<-chanelAppControl

	return
}

// Stop stop an application
func Stop() {
	chanelAppControl <- os.Interrupt
	<-chanelAppStop
}

// GetCmdArgs Инициализация параметров командной строки
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
