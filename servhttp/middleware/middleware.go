package middleware

import (
	"context"
	"crypto/rand"
	"io"
	"net/http"
	"time"

	"github.com/sungora/app/request"

	"github.com/sungora/app/core"
)

const (
	Num               = "0123456789"
	Strdown           = "abcdefghijklmnopqrstuvwxyz"
	Strup             = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Symbol            = "~!@#$%^&*_+-="
	KeyRW      string = "RW"
	KeySession string = "SESSION"
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

// NotFound обработчик не реализованных запросов
func NotFound(w http.ResponseWriter, r *http.Request) {
	rw := request.NewIn(w, r)
	rw.Static(core.Cfg.DirWww + r.URL.Path)
}

// newRandomString generates password key of a specified length (a-z0-9.)
func newRandomString(length int) string {
	return randChar(length, []byte(Strdown+Strup+Num))
}

func randChar(length int, chars []byte) string {
	pword := make([]byte, length)
	data := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, data); err != nil {
			panic(err)
		}
		for _, c := range data {
			if c >= maxrb {
				continue
			}
			pword[i] = chars[c%clen]
			i++
			if i == length {
				return string(pword)
			}
		}
	}
	panic("unreachable")
}
