package core

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"gopkg.in/sungora/app.v1/lg"
	"gopkg.in/sungora/app.v1/tool"
	"gopkg.in/sungora/app.v1/workflow"
)

type config struct {
	Main       configMain
	Mysql      configMysql
	Postgresql configPostgresql
	Log        lg.Config
	Workflow   workflow.Config
}

type configMain struct {
	TimeZone       string        // Временная зона
	DriverDB       string        // Драйвер DB
	SessionTimeout time.Duration // Время жизни сессии в секундах
	Host           string        // Хост веб сервера
	Port           int           // Порт веб сервера
	Mode           string        // Режим работы приложения
	Static         string        // Папка для статических данными (css, js, img, etc...)
	Template       string        // Папка для шаблонов
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

var (
	chanelAppStop    = make(chan os.Signal, 1) // Канал управления и остановкой приложения
	chanelAppControl = make(chan os.Signal, 1) // Канал управления и остановкой приложения
	Config           = new(config)
	DB               *gorm.DB
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
	if store, err = newHttp(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	defer store.Close()
	fmt.Fprintf(
		os.Stdout,
		"service start success: http://%s:%d\n",
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
	<-chanelAppStop
}

// newHTTP создание и запуск сервера
func newHttp() (store net.Listener, err error) {
	Server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", Config.Main.Host, Config.Main.Port),
		Handler:        new(httpHandler),
		ReadTimeout:    time.Second * time.Duration(300),
		WriteTimeout:   time.Second * time.Duration(300),
		MaxHeaderBytes: 1048576,
	}
	for i := 5; i > 0; i-- {
		store, err = net.Listen("tcp", Server.Addr)
		time.Sleep(time.Millisecond * 100)
		if err == nil {
			break
		}
	}
	if err == nil && store != nil {
		go Server.Serve(store)
		return
	} else if err == nil {
		return nil, errors.New("service start unknown error")
	}
	return nil, err
}

type httpHandler struct {
}

// ServeHTTP Точка входа запроса (в приложение).
func (handler *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		control ControllerFace
	)

	// search controller & static
	if control, err = Route.GetRoute(r.URL.Path); err != nil {
		rwH := new(rw)
		rwH.ResponseStatic(tool.DirWww + r.URL.Path)
		return
	}

	// initialization controller
	control.Init(w, r)

	// execute controller
	switch r.Method {
	case "GET":
		control.GET()
	case "POST":
		control.POST()
	case "PUT":
		control.PUT()
	case "DELETE":
		control.DELETE()
	case "OPTIONS":
		control.OPTIONS()
	default:
		rwH := new(rw)
		rwH.ResponseStatic(tool.DirTpl + "/404.html")
		return
	}

	// response controller
	control.Response()
}

// // search controller method
// objValue := reflect.ValueOf(control)
// met := objValue.MethodByName(r.Method)
// if met.IsValid() == false {
// 	control := &core.Controller{}
// 	control.Init(w, r, self.config)
// 	control.RW.ResponseHtml([]byte("page not found (m)"), 404)
// 	lg.Error("not found method [%s] of control", r.Method)
// 	return
// }
//
// // пример передачи параметров в метод
// var in = make([]reflect.Value, 0)
// var params []interface{}
// for i := range params {
// 	in = append(in, reflect.ValueOf(params[i]))
// }
//
// // execute method of controller
// out := met.Call(in)
// if nil != out[0].Interface() {
// 	lg.Error(out[0].Interface().(error))
// }

// GetCmdArgs Инициализация параметров командной строки
// func GetCmdArgs() (mode string, err error) {
// 	if len(os.Args) > 1 {
// 		mode = os.Args[1]
// 	}
// 	// - проверки
// 	if mode == `-h` || mode == `-help` || mode == `--help` {
// 		var str string
// 		str += "Usage of %s: %s [mode]\n"
// 		str += "\t mode: Режим запуска приложения\n"
// 		str += "\t\t install - Установка как сервиса в ОС\n"
// 		str += "\t\t uninstall - Удаление сервиса из ОС\n"
// 		str += "\t\t restart - Перезапуск ранее установленного сервиса\n"
// 		str += "\t\t start - Запуск ранее установленного сервиса\n"
// 		str += "\t\t stop - Остановка ранее установленного сервиса\n"
// 		str += "\t\t run - Прямой запуск (выход по 'Ctrl+C')\n"
// 		str += "\t\t если параметр опущен работает в режиме run\n"
// 		fmt.Fprintf(os.Stderr, str, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
// 		return "", errors.New("Help startup request")
// 	}
// 	return
// }
