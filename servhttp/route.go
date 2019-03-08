package servhttp

import (
	"net/http"

	"github.com/go-chi/chi"
)

// MountRoutes монтирование роутинга и его обработчиков подключаемыми модулями
func MountRoutes(pattern string, mount func() http.Handler) {
	route.Mount(pattern, mount())
}

// MiddlewareRoot set root middleware
func MiddlewareRoot(handler func(next http.Handler) http.Handler) {
	route.Use(handler)
}

// NotFound set page NotFound
func NotFound(handlerFn http.HandlerFunc) {
	route.NotFound(handlerFn)
}

// GetRoute get handler route
func GetRoute() *chi.Mux {
	return route
}
