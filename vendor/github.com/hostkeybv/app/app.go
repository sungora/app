package app

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/hostkeybv/app/tool"
)

type Config struct {
	Mode           string        // Режим работы приложения
	Proto          string        // Режим работы приложения
	TimeZone       string        // Временная зона
	DriverDB       string        // Драйвер DB
	SessionTimeout time.Duration // Время жизни сессии в секундах
	Host           string        // Хост веб сервера
	Port           int           // Порт веб сервера
}

var config Config
var store net.Listener

func Init(c Config) (err error) {
	config = c
	// сервер
	if config.Proto == "http" {
		http.HandleFunc("/", httpHandler) // установим роутер
		store, err = net.Listen("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
		if err != nil {
			return
		}
		go http.Serve(store, nil)
		fmt.Fprintf(
			os.Stdout,
			"service start success: %s:%d\n",
			c.Host,
			c.Port,
		)
	}
	// session
	if 0 < config.SessionTimeout {
		sessionGC()
	}
	return
}

func Wait() {
	if store != nil {
		store.Close()
	}
}

// httpHandler Точка входа запроса (в приложение).
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
