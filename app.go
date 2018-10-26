// Управление запуском и остановкой приложения
package app

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	"gopkg.in/sungora/app.v1/conf"
	"gopkg.in/sungora/app.v1/core"
	"gopkg.in/sungora/app.v1/lg"
	"gopkg.in/sungora/app.v1/workflow"
)

type Config struct {
	Main       ConfigMain
	Mysql      ConfigMysql
	Postgresql ConfigPostgresql
	Log        lg.Config
	Workflow   workflow.Config
}

type ConfigMain struct {
	TimeZone    string
	DriverDB    string // Драйвер DB
	Mode        string // Режим работы приложения
	Host        string
	Port        int
}

type ConfigMysql struct {
	Host     string // протокол, хост и порт подключения
	Name     string // Имя базы данных
	Login    string // Логин к базе данных
	Password string // Пароль к базе данных
	Charset  string // Кодировка данных (utf-8 - по умолчанию)
}

type ConfigPostgresql struct {
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
func Start(fileConfigName string) (code int) {
	defer func() { // контроль завершение работы приложения
		chanelAppStop <- os.Interrupt
	}()
	var (
		err   error
		store net.Listener
	)

	// config
	var configApp *Config
	if _, err = toml.DecodeFile(conf.ConfigDir+fileConfigName+".toml", &configApp); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	// Инициализация временной зоны
	if loc, err := time.LoadLocation(configApp.Main.TimeZone); err == nil {
		conf.TimeLocation = loc
	} else {
		conf.TimeLocation = time.UTC
	}

	// logs
	if err = lg.Start(configApp.Log); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	defer lg.Wait()

	// base controller
	switch configApp.Main.DriverDB {
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

	// workflow
	if configApp.Workflow.IsWorkflow == true {
		if err = workflow.Start(configApp.Workflow); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
		defer workflow.Wait()
	}

	// web server - application
	if store, err = newWeb(configApp.Main); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	defer store.Close()
	fmt.Fprintln(os.Stdout, "web app start success")

	// The correctness of the application is closed by a signal
	signal.Notify(chanelAppControl, os.Interrupt)
	<-chanelAppControl

	return
}

// stop Stop an application
func Stop() {
	chanelAppControl <- os.Interrupt
	<-chanelAppStop
}
