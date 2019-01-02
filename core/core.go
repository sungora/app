package core

import (
	"fmt"
	"gopkg.in/sungora/app.v1/lg"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"gopkg.in/sungora/app.v1/tool"
	"gopkg.in/sungora/app.v1/uploader"
	"gopkg.in/sungora/app.v1/workflow"
)

var (
	chanelAppControl = make(chan os.Signal, 1) // Канал управления остановкой приложения
	Config           = new(config)
	DB               *gorm.DB
)

// Start Launch an application
func Start() (code int) {
	defer func() { // контроль завершение работы приложения
		chanelAppControl <- os.Interrupt
	}()
	var (
		err   error
		store net.Listener
	)

	// config
	path := tool.DirConfig + string(os.PathSeparator) + "main.toml"
	if _, err = toml.DecodeFile(path, Config); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	Config.Main.SessionTimeout *= time.Second

	// Инициализация временной зоны
	if loc, err := time.LoadLocation(Config.Main.TimeZone); err == nil {
		tool.TimeLocation = loc
	} else {
		tool.TimeLocation = time.UTC
	}

	// logs
	if err = lg.Start(Config.Log); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	defer lg.Wait()

	// Модуль загрузки файлов и получение их по идентификатору
	if err = uploader.Init(tool.DirWww, 30); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	// core
	switch Config.Main.DriverDB {
	case "mysql":
		if DB, err = gorm.Open("mysql", fmt.Sprintf(
			"%s:%s@%s/%s?charset=%s&parseTime=True&loc=Local&timeout=3s",
			Config.Mysql.Login,
			Config.Mysql.Password,
			Config.Mysql.Host,
			Config.Mysql.Name,
			Config.Mysql.Charset,
		)); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
		defer DB.Close()
	case "postgresql":
		if DB, err = gorm.Open("postgres", fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s",
			Config.Postgresql.Host,
			Config.Postgresql.Port,
			Config.Postgresql.Login,
			Config.Postgresql.Name,
			Config.Postgresql.Password,
		)); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
		defer DB.Close()
	}

	// session
	if 0 < Config.Main.SessionTimeout {
		sessionGC()
	}

	// workflow
	if Config.Workflow.IsWorkflow == true {
		if err = workflow.Start(Config.Workflow); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
		defer workflow.Wait()
	}

	// service - application
	http.HandleFunc("/", httpHandler) // установим роутер
	store, err = net.Listen("tcp", fmt.Sprintf("%s:%d", Config.Main.Host, Config.Main.Port))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	go http.Serve(store, nil)
	defer store.Close()
	fmt.Fprintf(
		os.Stdout,
		"service start success: %s:%d\n",
		Config.Main.Host,
		Config.Main.Port,
	)
	// The correctness of the application is closed by a signal
	signal.Notify(chanelAppControl, os.Interrupt)
	<-chanelAppControl
	return
}

// Wait an application
func Wait() {
	chanelAppControl <- os.Interrupt
	<-chanelAppControl
}

// ServeHTTP Точка входа запроса (в приложение).
// func (handler *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
func httpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// search controller
	var control, err = Route.Get(r.URL.Path)
	// response static
	if err != nil {
		rwH := newRW(r, w)
		rwH.ResponseStatic(tool.DirWww + r.URL.Path)
		return
	}
	// initialization controller
	control.Init(w, r)
	// execute controller
	switch r.Method {
	case http.MethodGet:
		control.GET()
	case http.MethodPost:
		control.POST()
	case http.MethodPut:
		control.PUT()
	case http.MethodDelete:
		control.DELETE()
	case http.MethodOptions:
		control.OPTIONS()
	}
	// default response controller
	control.Response()
}
