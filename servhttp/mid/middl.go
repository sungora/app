package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/sungora/app/core"
)

const KeyRW string = "RW"

// Main (middleware)
// инициализация таймаута контекста для контроля времени выполениня запроса
// инициализация инструмента для ответа
func Main(d time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, KeyRW, core.NewRW(w, r))))
		})
	}
}

// NotFound
func NotFound(w http.ResponseWriter, r *http.Request) {
	var rw = r.Context().Value(KeyRW).(*core.RW)
	rw.ResponseStatic(core.Config.DirWww + r.URL.Path)
}
