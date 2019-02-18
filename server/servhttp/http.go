package servhttp

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
		if err = rw.ResponseStatic(core.Config.DirWww + r.URL.Path); err != nil {
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
