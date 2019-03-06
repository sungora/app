package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"

	"github.com/sungora/app/core"
	"github.com/sungora/app/lg"
	"github.com/sungora/app/servhttp"
	"github.com/sungora/app/servhttp/middleware"
)

func main() {




	// инициализация компонентов
	if 1 == core.Init() {
		os.Exit(1)
	}

	servhttp.RootMiddleware(middleware.Main(time.Second*5 - time.Millisecond))
	servhttp.NotFound(middleware.NotFound)
	servhttp.MountRoutes("/", Routes)

	// log
	lg.SetMessages(map[int]string{
		1000: "Message format Fmt from 1000",
		1001: "Message format Fmt from 1001",
		1002: "Message format Fmt from 1002",
		1003: "Message format Fmt from 1003",
		1004: "Message format Fmt from 1004",
		1005: "Message format Fmt from 1005",
	})

	// запуск приложения
	os.Exit(core.Start())
}

func Routes() http.Handler {
	r := chi.NewRouter()
	r.HandleFunc("/", Main)
	return r
}

// Main главная страница
func Main(w http.ResponseWriter, r *http.Request) {
	var rw = r.Context().Value(middleware.KeyRW).(*core.RW)
	rw.ResponseHtml("popcorn", 200)
	return

}
