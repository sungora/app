package server

import "net/http"

type serverHttp struct {
}

// ServeHTTP Точка входа запроса (в приложение).
func (server *serverHttp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Write([]byte(r.URL.Path))

	w.WriteHeader(http.StatusNotFound)
}
