package middles

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/cors"

	"github.com/sungora/app"
	"github.com/sungora/app/servhttp"
	"github.com/sungora/app/session"
	"github.com/sungora/app/request"
)

const (
	KeySession = "SESSION"
)

// TimeoutContext (middleware)
// инициализация таймаута контекста для контроля времени выполениня запроса
func TimeoutContext(d time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			d = time.Second*d - time.Millisecond
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
		token, _ := rw.CookieGet(app.Cfg.ServiceName)
		if token == "" {
			token = newRandomString(10)
			rw.CookieSet(app.Cfg.ServiceName, token)
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, KeySession, session.GetSession(token))
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
	rw.Static(app.Cfg.DirWork + r.URL.Path)
}
