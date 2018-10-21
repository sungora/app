package web

import (
	"errors"
)

var routes = make(map[string]ControllerFace)

func SetRouter(uri string, control ControllerFace) {
	routes[uri] = control
}
func GetRouter(uri string) (control ControllerFace, err error) {
	if _, ok := routes[uri]; ok {
		return routes[uri], nil
	}
	return nil, errors.New("route not found: " + uri)

}
