package middelw

import (
	"context"
	"github.com/sungora/app/servhttp"
	"net/http"
	"time"

	"github.com/go-chi/cors"

	"github.com/sungora/app/core"
	"github.com/sungora/app/request"
)

const (
	KeyRW      = "RW"
	KeySession = "SESSION"
)

// TimeoutContext (middleware)
// инициализация таймаута контекста для контроля времени выполениня запроса
func TimeoutContext(d time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Session Инициализация сессии
func Session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := request.NewIn(w, r)
		token, _ := rw.CookieGet(core.Cfg.ServiceName)
		if token == "" {
			token = newRandomString(10)
			rw.CookieSet(core.Cfg.ServiceName, token)
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, KeySession, core.GetSession(token))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Cors добавление заголовка Cors
func Cors(cfg servhttp.Cors) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   cfg.AllowedMethods,
		AllowedHeaders:   cfg.AllowedHeaders,
		ExposedHeaders:   cfg.ExposedHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           cfg.MaxAge, // Maximum value not ignored by any of major browsers
	})
}

// NotFound обработчик не реализованных запросов
func NotFound(w http.ResponseWriter, r *http.Request) {
	rw := request.NewIn(w, r)
	rw.Static(core.Cfg.DirWww + r.URL.Path)
}
