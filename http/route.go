package http

import (
	"net/http"
)

// MountRoutes монитрование роутинга и его обработчиков подключаемыми модулями
func MountRoutes(pattern string, mount func() http.Handler) {
	route.Mount(pattern, mount())
}
