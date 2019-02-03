package server

import (
	"net/http"

	"github.com/sungora/app/core"
	"github.com/sungora/app/lg"
)

type serverHttp struct {
}

// ServeHTTP Точка входа запроса (в приложение).
func (server *serverHttp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// search controller
	var control, err = core.Route.Get(r.URL.Path)
	// response static
	if err != nil {
		var rw = core.NewRW(r, w)
		if err = rw.ResponseStatic(core.DirWww + r.URL.Path); err != nil {
			lg.Error(err)
		}
		return
	}
	// initialization controller
	if err = control.Init(w, r); err != nil {
		lg.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// execute controller
	switch r.Method {
	case http.MethodGet:
		err = control.GET()
	case http.MethodPost:
		err = control.POST()
	case http.MethodPut:
		err = control.PUT()
	case http.MethodDelete:
		err = control.DELETE()
	case http.MethodOptions:
		err = control.OPTIONS()
	}
	if err != nil {
		lg.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// default response controller
	if err = control.Response(); err != nil {
		lg.Error(err)
	}
}
