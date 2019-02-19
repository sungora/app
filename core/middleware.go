package core

import (
	"bytes"
	"context"
)

type middleware struct {
}

func MidlSession(ctx context.Context, rw *RW) (context.Context, *RW) {
	// initialization session
	var token string
	var err error
	if 0 < Config.SessionTimeout {
		if token, err = rw.CookieGet(Config.ServiceName); err != nil {
			return ctx, rw
		}
		if token == "" {
			token = NewRandomString(10)
			rw.CookieSet(Config.ServiceName, token)
		}
		ctx = context.WithValue(ctx, "Session", GetSession(token))
	}
	return ctx, rw
}

func MidlResponseDefault(ctx context.Context, rw *RW) (context.Context, *RW) {
	if rw.isResponse {
		return ctx, rw
	}
	var (
		err           error
		Functions     = make(map[string]interface{})
		Variables     = make(map[string]interface{})
		TplLayout     = Config.DirWww + "/layout/index.html"
		TplController = Config.DirWww + "/controllers"
	)

	// шаблон контроллера
	var buf bytes.Buffer
	if buf, err = TplCompilation(TplController, Functions, Variables); err != nil {
		rw.ResponseHtml("<H1>Internal Server Error</H1>", 500)
		return ctx, rw
	}
	if TplLayout == "" {
		rw.ResponseHtml(buf.String(), 200)
		return ctx, rw
	}
	// шаблон макета
	Variables["Content"] = buf.String()
	if buf, err = TplCompilation(TplLayout, Functions, Variables); err != nil {
		rw.ResponseHtml("<H1>Internal Server Error</H1>", 500)
		return ctx, rw
	}
	rw.ResponseHtml(buf.String(), 200)
	return ctx, rw
}
