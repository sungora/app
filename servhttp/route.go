package servhttp

import (
	"net/http"
)

// MountRoutes монтирование роутинга и его обработчиков подключаемыми модулями
func MountRoutes(pattern string, mount func() http.Handler) {
	route.Mount(pattern, mount())
}

// RootMiddleware set root middleware
func RootMiddleware(handler func(next http.Handler) http.Handler) {
	route.Use(handler)
}

// NotFound set page NotFound
func NotFound(handlerFn http.HandlerFunc) {
	route.NotFound(handlerFn)
}
